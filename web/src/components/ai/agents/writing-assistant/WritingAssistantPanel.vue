<template>
  <div class="writing-assistant">
    <div v-if="!aiEnabled" class="ai-disabled-tip">
      <el-empty description="AI 功能不可用" :image-size="72" />
      <div class="tip-text">{{ disabledReason }}</div>
      <el-button size="small" class="retry-btn" @click="checkStatus">重新检测</el-button>
    </div>

    <template v-else>
      <!-- 快捷动作 -->
      <div class="action-grid">
        <el-button
          v-for="action in quickActions"
          :key="action.key"
          size="small"
          :disabled="action.disabled.value || streaming"
          @click="runAction(action)"
        >
          {{ action.label }}
        </el-button>
      </div>

      <!-- 自定义指令 -->
      <div class="custom-input">
        <el-input
          v-model="instruction"
          type="textarea"
          :rows="2"
          placeholder="输入自定义指令，如：把这段改成问答体"
          :disabled="streaming"
        />
        <el-button
          type="primary"
          size="small"
          class="send-btn"
          :disabled="!instruction.trim() || streaming"
          @click="runCustom"
        >
          发送
        </el-button>
      </div>

      <!-- 结果区 -->
      <div v-if="streaming || resultText || diffBlocks.length" class="result-area">
        <div class="result-header">
          <span class="result-title">{{ resultTitle }}</span>
          <el-button
            v-if="streaming"
            size="small"
            type="danger"
            plain
            @click="abortStream"
          >
            停止生成
          </el-button>
        </div>

        <!-- diff 模式 -->
        <template v-if="mode === 'diff'">
          <ParagraphDiff :blocks="diffBlocks" @change="refreshDiffStats" />
          <div class="result-actions">
            <el-button type="primary" size="small" @click="applyDiff">
              应用修改
            </el-button>
            <el-button size="small" @click="reset">丢弃</el-button>
          </div>
        </template>

        <!-- 普通结果模式 -->
        <template v-else>
          <div v-if="resultText" class="result-body" v-html="renderedResult" />
          <div v-else-if="streaming" class="streaming-hint">
            {{ streamingHint }}<span class="cursor">▌</span>
          </div>
          <div v-if="resultActions.length" class="result-actions">
            <el-button
              v-for="act in resultActions"
              :key="act.key"
              type="primary"
              size="small"
              :disabled="streaming"
              @click="act.handler"
            >
              {{ act.label }}
            </el-button>
            <el-button size="small" :disabled="streaming" @click="reset">丢弃</el-button>
          </div>
        </template>
      </div>

      <!-- 历史轮次 -->
      <div v-if="history.length" class="history-bar">
        <span>对话上下文 {{ history.length }} 轮</span>
        <el-button link size="small" @click="history = []">清空</el-button>
      </div>
    </template>
  </div>
</template>

<script setup>
  import { computed, ref, watch } from 'vue'
  import { Marked } from 'marked'
  import { markedHighlight } from 'marked-highlight'
  import hljs from 'highlight.js'
  import { ElMessage } from 'element-plus'
  import { useAiStore } from '@/pinia/modules/ai'
  import { streamAiChat, getAiStatus, generateSummary, suggestTags } from '@/api/blog/ai'
  import ParagraphDiff from './ParagraphDiff.vue'
  import { diffMarkdownBlocks, applyDiffBlocks } from './diff'

  const aiStore = useAiStore()

  const marked = new Marked(
    { gfm: true, breaks: true },
    markedHighlight({
      langPrefix: 'hljs language-',
      highlight(code, lang) {
        const language = hljs.getLanguage(lang) ? lang : 'plaintext'
        return hljs.highlight(code, { language }).value
      }
    })
  )

  // ---- 状态 ----
  const aiEnabled = ref(true)
  const disabledReason = ref('')
  const streaming = ref(false)
  const mode = ref('idle') // idle | result | diff
  const instruction = ref('')
  const resultText = ref('')
  const resultTitle = ref('')
  const streamingHint = ref('')
  const diffBlocks = ref([])
  const history = ref([])
  let currentAction = null
  let currentSource = '' // 润色/改写的原始选区
  let streamHandle = null

  const editorCtx = computed(() => aiStore.contexts.editor)
  const hasEditor = computed(() => Boolean(editorCtx.value))
  const hasSelection = computed(() => {
    try {
      return Boolean(editorCtx.value?.getSelection?.()?.text)
    } catch (_) {
      return false
    }
  })
  const hasContent = computed(() => {
    try {
      return Boolean(editorCtx.value?.getFullText?.())
    } catch (_) {
      return false
    }
  })

  checkStatus()
  async function checkStatus() {
    disabledReason.value = ''
    try {
      const res = await getAiStatus()
      if (res.data?.enabled) {
        aiEnabled.value = true
      } else {
        aiEnabled.value = false
        disabledReason.value = '未配置默认模型：请在「AI 模型配置」中将一个启用中的模型设为默认'
      }
    } catch (error) {
      aiEnabled.value = false
      const status = error?.response?.status
      if (status === 403) {
        disabledReason.value =
          '没有接口权限（403）：请在 角色管理 → 分配API权限 中，为当前角色勾选 /blog/ai/* 四个接口'
      } else if (status === 401) {
        disabledReason.value = '登录已过期，请重新登录后再试'
      } else if (!error?.response) {
        disabledReason.value = '无法连接服务器，请确认后端已启动'
      } else {
        disabledReason.value = error?.response?.data?.msg || '探测 AI 状态失败'
      }
    }
  }

  // 每次打开 AI 抽屉都重新探测（权限/配置变更后无需刷新页面）
  watch(
    () => aiStore.dockVisible,
    (visible) => {
      if (visible && !aiEnabled.value) checkStatus()
    }
  )

  // ---- 快捷动作定义 ----
  const quickActions = [
    {
      key: 'polish', label: '润色',
      disabled: computed(() => !hasSelection.value || !hasEditor.value),
      hint: '请先在编辑器中选中要润色的文本'
    },
    {
      key: 'rewrite', label: '改写',
      disabled: computed(() => !hasSelection.value || !hasEditor.value),
      hint: '请先在编辑器中选中要改写的文本'
    },
    {
      key: 'continue', label: '续写',
      disabled: computed(() => !hasContent.value || !hasEditor.value),
      hint: '请先在编辑器中输入一些正文'
    },
    {
      key: 'outline', label: '生成大纲',
      disabled: computed(() => false),
      hint: ''
    },
    {
      key: 'title', label: '起标题',
      disabled: computed(() => !hasContent.value),
      hint: '请先输入正文或标题'
    },
    {
      key: 'summary', label: '生成摘要',
      disabled: computed(() => !hasContent.value),
      hint: '请先输入正文'
    },
    {
      key: 'suggest-tags', label: '推荐标签',
      disabled: computed(() => !hasContent.value),
      hint: '请先输入正文'
    }
  ]

  const renderedResult = computed(() => marked.parse(resultText.value || ''))

  const resultActions = computed(() => {
    const actions = []
    if (currentAction === 'summary') {
      actions.push({ key: 'fill', label: '回填摘要', handler: fillDescription })
    }
    if (currentAction === 'suggest-tags') {
      actions.push({ key: 'apply', label: '回填标签', handler: applyTags })
    }
    if (hasEditor.value && (currentAction === 'continue' || currentAction === 'outline')) {
      actions.push({ key: 'insert', label: '插入到光标处', handler: insertResult })
    }
    actions.push({ key: 'copy', label: '复制', handler: copyResult })
    return actions
  })

  // ---- 执行 ----
  const runAction = (action) => {
    if (action.disabled.value) {
      if (action.hint) ElMessage.warning(action.hint)
      return
    }
    runActionByKey(action.key)
  }

  const runActionByKey = (key) => {
    if (key === 'summary') return runSummary()
    if (key === 'suggest-tags') return runSuggestTags()
    return runChatAction(key)
  }

  const buildPayload = (action, extra = {}) => {
    const editor = editorCtx.value
    const payload = {
      action,
      title: editor?.getTitle?.() || '',
      content: editor?.getFullText?.() || '',
      ...extra
    }
    const selection = editor?.getSelection?.()
    payload.selection = selection?.text || ''
    payload.cursorContext = editor?.getCursorContext?.() || payload.content
    return payload
  }

  const runChatAction = async (action) => {
    const selection = editorCtx.value?.getSelection?.()
    currentSource = selection?.text || ''
    currentAction = action
    mode.value = 'result'
    resultText.value = ''
    diffBlocks.value = []
    resultTitle.value = {
      polish: '润色结果', rewrite: '改写结果', continue: '续写结果',
      outline: '生成大纲', title: '标题建议', custom: 'AI 结果'
    }[action] || 'AI 结果'
    streamingHint.value = '正在思考'
    await startStream(buildPayload(action))
  }

  const runCustom = async () => {
    if (!instruction.value.trim()) return
    currentAction = 'custom'
    mode.value = 'result'
    resultText.value = ''
    diffBlocks.value = []
    resultTitle.value = 'AI 结果'
    streamingHint.value = '正在思考'
    await startStream(buildPayload('custom', { instruction: instruction.value.trim() }))
  }

  const runSummary = async () => {
    currentAction = 'summary'
    mode.value = 'result'
    resultText.value = ''
    diffBlocks.value = []
    resultTitle.value = '文章摘要'
    streaming.value = true
    try {
      const res = await generateSummary(buildPayload('summary'))
      resultText.value = res.data?.summary || ''
    } catch (error) {
      ElMessage.error(error?.msg || '生成摘要失败')
    } finally {
      streaming.value = false
    }
  }

  const runSuggestTags = async () => {
    currentAction = 'suggest-tags'
    mode.value = 'result'
    resultText.value = ''
    diffBlocks.value = []
    resultTitle.value = '分类与标签建议'
    streaming.value = true
    try {
      const res = await suggestTags(buildPayload('suggest-tags'))
      const data = res.data || {}
      resultText.value = [
        data.category ? `分类：${data.category}` : '',
        data.tags?.length ? `标签：${data.tags.join('、')}` : '',
        data.newTags?.length ? `建议新建：${data.newTags.join('、')}` : ''
      ].filter(Boolean).join('\n') || '没有建议'
      tagSuggestion.value = data
    } catch (error) {
      ElMessage.error(error?.msg || '推荐标签失败')
    } finally {
      streaming.value = false
    }
  }

  let tagSuggestion = ref(null)

  const startStream = (payload) => {
    streaming.value = true
    const userContent = payload.instruction || payload.selection || payload.content || ''
    streamHandle = streamAiChat(
      { ...payload, history: history.value.slice(-6) },
      {
        onDelta: (delta) => {
          resultText.value += delta
        },
        onTool: (tool) => {
          if (tool.name === 'search_my_blogs') streamingHint.value = '正在检索你的历史文章'
          else if (tool.name === 'get_blog_content') streamingHint.value = '正在阅读相关文章'
          else streamingHint.value = '正在调用工具'
        },
        onError: (message) => {
          ElMessage.error(message)
        }
      }
    )

    return streamHandle.promise.then(({ aborted }) => {
      streaming.value = false
      streamHandle = null
      if (aborted) return
      if (!resultText.value) return

      history.value.push(
        { role: 'user', content: userContent.slice(0, 2000) },
        { role: 'assistant', content: resultText.value.slice(0, 2000) }
      )
      if (history.value.length > 6) history.value = history.value.slice(-6)

      // 润色/改写且原文来自选区 → 进入 diff 模式
      if ((payload.action === 'polish' || payload.action === 'rewrite') && currentSource) {
        diffBlocks.value = diffMarkdownBlocks(currentSource, resultText.value)
        mode.value = 'diff'
      }
    })
  }

  const abortStream = () => {
    streamHandle?.abort()
    streaming.value = false
    streamHandle = null
  }

  // ---- 结果操作 ----
  const insertResult = () => {
    editorCtx.value?.insertAtCursor?.(resultText.value)
    ElMessage.success('已插入到光标处')
    reset()
  }

  const fillDescription = () => {
    if (editorCtx.value?.fillDescription?.(resultText.value)) {
      ElMessage.success('已回填到文章摘要')
      reset()
    } else {
      ElMessage.warning('当前页面不支持回填，请手动复制')
    }
  }

  const applyTags = () => {
    if (!tagSuggestion.value) return
    if (editorCtx.value?.applySuggestion?.(tagSuggestion.value)) {
      ElMessage.success('已回填分类与标签')
      reset()
    } else {
      ElMessage.warning('当前页面不支持回填，请手动复制')
    }
  }

  const refreshDiffStats = () => {
    /* ParagraphDiff 内部自行统计，这里仅作为 change 钩子 */
  }

  const applyDiff = () => {
    const finalText = applyDiffBlocks(diffBlocks.value)
    editorCtx.value?.replaceSelection?.(finalText)
    ElMessage.success('已应用修改')
    reset()
  }

  const copyResult = async () => {
    try {
      await navigator.clipboard.writeText(resultText.value)
      ElMessage.success('已复制')
    } catch (_) {
      ElMessage.error('复制失败')
    }
  }

  const reset = () => {
    mode.value = 'idle'
    resultText.value = ''
    diffBlocks.value = []
    currentAction = null
    currentSource = ''
    tagSuggestion.value = null
  }
</script>

<style scoped lang="scss">
.writing-assistant {
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  overflow: hidden;
}

.ai-disabled-tip {
  padding: 24px 0;

  .tip-text {
    margin-top: -12px;
    padding: 0 16px;
    color: #909399;
    font-size: 13px;
    line-height: 1.7;
    text-align: center;
  }

  .retry-btn {
    display: block;
    margin: 12px auto 0;
  }
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;

  .el-button {
    width: 100%;
    margin: 0;
  }
}

.custom-input {
  display: flex;
  flex-direction: column;
  gap: 8px;

  .send-btn {
    align-self: flex-end;
  }
}

.result-area {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  padding: 10px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #fafbfc;
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;

  .result-title {
    font-weight: 600;
    font-size: 14px;
  }
}

.result-body {
  flex: 1;
  min-height: 0;
  margin-bottom: 8px;
  overflow: auto;
  padding: 8px;
  border-radius: 4px;
  background: #fff;
  font-size: 13px;
  line-height: 1.7;

  :deep(pre) {
    overflow: auto;
    padding: 10px;
    border-radius: 4px;
    background: #f6f8fa;
  }

  :deep(code) {
    font-family: ui-monospace, Menlo, Consolas, monospace;
    font-size: 12px;
  }
}

.streaming-hint {
  padding: 8px;
  color: #909399;
  font-size: 13px;

  .cursor {
    animation: blink 1s infinite;
  }
}

@keyframes blink {
  50% {
    opacity: 0;
  }
}

.result-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: auto;
  padding-top: 8px;
}

.history-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #909399;
  font-size: 12px;
}
</style>
