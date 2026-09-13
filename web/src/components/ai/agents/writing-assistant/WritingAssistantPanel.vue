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
          :disabled="action.disabled.value || streaming || editorDiffOpen"
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
          :disabled="streaming || editorDiffOpen"
        />
        <el-button
          type="primary"
          size="small"
          class="send-btn"
          :disabled="!instruction.trim() || streaming || editorDiffOpen"
          @click="runCustom"
        >
          发送
        </el-button>
      </div>

      <!-- 结果区 -->
      <div v-if="streaming || resultText || editorDiffOpen" class="result-area">
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

        <!-- 润色/改写完成：diff 已铺在正文编辑器中 -->
        <div v-if="editorDiffOpen" class="editor-diff-notice">
          <div>已在正文编辑器中打开 diff 对比：逐块选择「采用 AI 版 / 保留原文」，右侧预览可实时查看采纳效果。</div>
          <div class="notice-actions">
            <el-button size="small" @click="discardEditorDiff">放弃对比</el-button>
          </div>
        </div>

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
        <span>对话上下文 {{ Math.floor(history.length / 2) }} 轮</span>
        <el-button link size="small" @click="history = []">清空</el-button>
      </div>
    </template>
  </div>
</template>

<script setup>
  import { computed, onBeforeUnmount, ref, watch } from 'vue'
  import { renderSafeMarkdown } from '@/utils/safeMarkdown'
  import { ElMessage } from 'element-plus'
  import { useAiStore } from '@/pinia/modules/ai'
  import { streamAiChat, getAiStatus, generateSummary, suggestTags } from '@/api/blog/ai'

  const aiStore = useAiStore()

  // ---- 状态 ----
  const aiEnabled = ref(true)
  const disabledReason = ref('')
  const streaming = ref(false)
  const instruction = ref('')
  const resultText = ref('')
  const resultTitle = ref('')
  const streamingHint = ref('')
  const editorDiffOpen = computed(() => aiStore.diff.active)
  const history = ref([])
  const currentAction = ref(null)
  let streamHandle = null
  const resultReady = ref(false)
  const tagSuggestion = ref(null)
  let resultOwner = null

  const editorCtx = computed(() => aiStore.contexts.editor)
  const hasEditor = computed(() => Boolean(editorCtx.value))

  // textarea 选区不是响应式数据，监听原生事件刷新，保证"润色/改写"按钮状态随选中实时变化
  const selectionText = ref('')
  const syncSelection = () => {
    try {
      selectionText.value = editorCtx.value?.getSelection?.()?.text || ''
    } catch (_) {
      selectionText.value = ''
    }
  }
  document.addEventListener('selectionchange', syncSelection)
  onBeforeUnmount(() => {
    document.removeEventListener('selectionchange', syncSelection)
    abortStream()
  })

  const hasSelection = computed(() => Boolean(selectionText.value))
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
      if (res.data?.quotaError) {
        aiEnabled.value = false
        disabledReason.value = res.data.quotaError
      } else if (res.data?.enabled) {
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

  // 每次打开 AI 抽屉都重新探测（权限/配置变更后无需刷新页面），并同步当前选区
  watch(
    () => aiStore.dockVisible,
    (visible) => {
      if (!visible) return
      syncSelection()
      if (!aiEnabled.value) checkStatus()
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

  const renderedResult = computed(() => renderSafeMarkdown(resultText.value))

  const resultActions = computed(() => {
    if (!resultText.value) return []
    if (!resultReady.value) return [{ key: 'copy', label: '复制', handler: copyResult }]
    const actions = []
    if (currentAction.value === 'summary') {
      actions.push({ key: 'fill', label: '回填摘要', handler: fillDescription })
    }
    if (currentAction.value === 'suggest-tags') {
      actions.push({ key: 'apply', label: '回填标签', handler: applyTags })
    }
    if (hasEditor.value && (currentAction.value === 'continue' || currentAction.value === 'outline')) {
      actions.push({ key: 'insert', label: '插入到光标处', handler: insertResult })
    }
    actions.push({ key: 'copy', label: '复制', handler: copyResult })
    return actions
  })

  // ---- 执行 ----
  const runAction = (action) => {
    if (streaming.value || editorDiffOpen.value) return
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
    const needsSelection = action === 'polish' || action === 'rewrite'
    const snapshot = needsSelection ? editorCtx.value?.captureSelection?.() : null
    if (needsSelection && !snapshot) return ElMessage.warning('请先在当前文章中选择要修改的文本')
    currentAction.value = action
    resultText.value = ''
    resultTitle.value = {
      polish: '润色结果', rewrite: '改写结果', continue: '续写结果',
      outline: '生成大纲', title: '标题建议', custom: 'AI 结果'
    }[action] || 'AI 结果'
    streamingHint.value = '正在思考'
    const payload = buildPayload(action)
    if (snapshot) {
      payload.selection = snapshot.text
      payload.content = snapshot.content
    }
    await startStream(payload, snapshot)
  }

  const runCustom = async () => {
    if (!instruction.value.trim() || streaming.value || editorDiffOpen.value) return
    currentAction.value = 'custom'
    resultText.value = ''
    resultTitle.value = 'AI 结果'
    streamingHint.value = '正在思考'
    await startStream(buildPayload('custom', { instruction: instruction.value.trim() }))
  }

  const runSummary = () => runSingle('summary', generateSummary)
  const runSuggestTags = () => runSingle('suggest-tags', suggestTags)

  const runSingle = async (action, request) => {
    const handle = new AbortController()
    beginTask(handle)
    currentAction.value = action
    resultText.value = ''
    resultTitle.value = action === 'summary' ? '文章摘要' : '分类与标签建议'
    streamingHint.value = '正在思考'
    try {
      const res = await request(buildPayload(action), handle.signal)
      if (streamHandle !== handle) return
      if (action === 'summary') {
        const summary = res.data?.summary
        if (typeof summary !== 'string' || !summary.trim()) throw new Error('AI 未返回有效摘要')
        resultText.value = summary
      } else {
        const data = res.data
        if (!data || !Array.isArray(data.tags) || !Array.isArray(data.newTags) ||
          [...data.tags, ...data.newTags].some(value => typeof value !== 'string') ||
          (data.category != null && typeof data.category !== 'string')) throw new Error('AI 返回了无效的标签建议')
        resultText.value = [
          data.category ? `分类：${data.category}` : '',
          data.tags.length ? `标签：${data.tags.join('、')}` : '',
          data.newTags.length ? `建议新建：${data.newTags.join('、')}` : ''
        ].filter(Boolean).join('\n') || '没有建议'
        tagSuggestion.value = data
      }
      resultReady.value = true
    } catch (error) {
      if (streamHandle === handle && !handle.signal.aborted) ElMessage.error(error?.response?.data?.msg || error?.message || 'AI 请求失败')
    } finally {
      if (streamHandle === handle) {
        streaming.value = false
        streamHandle = null
      }
    }
  }

  const captureOwner = () => {
    const context = editorCtx.value
    const state = context?.getEditorState?.()
    return { context, editorId: state?.editorId, documentId: state?.documentId }
  }
  const ownsResult = () => {
    const owner = captureOwner()
    return resultOwner && owner.context === resultOwner.context &&
      owner.editorId === resultOwner.editorId && owner.documentId === resultOwner.documentId
  }
  const beginTask = (handle) => {
    streamHandle?.abort()
    streamHandle = handle
    resultReady.value = false
    tagSuggestion.value = null
    resultOwner = captureOwner()
    streaming.value = true
  }

  const startStream = (payload, snapshot = null) => {
    resultReady.value = false
    const userContent = payload.instruction || payload.selection || payload.content || ''
    let generatedText = ''
    let completed = false
    const handle = streamAiChat(
      { ...payload, history: history.value.slice(-6) },
      {
        onDelta: (delta) => {
          if (streamHandle !== handle) return
          generatedText += delta
          resultText.value = generatedText
        },
        onTool: (tool) => {
          if (streamHandle !== handle) return
          if (tool.name === 'search_my_blogs') streamingHint.value = '正在检索你的历史文章'
          else if (tool.name === 'get_blog_content') streamingHint.value = '正在阅读相关文章'
          else streamingHint.value = '正在调用工具'
        },
        onError: (message) => {
          if (streamHandle !== handle) return
          ElMessage.error(message)
        },
        onDone: () => {
          completed = true
        }
      }
    )

    beginTask(handle)
    return handle.promise.then(({ aborted, errorMessage }) => {
      if (streamHandle !== handle) return
      streaming.value = false
      streamHandle = null
      if (aborted || errorMessage) return
      if (!generatedText.trim()) return ElMessage.warning('AI 未返回有效内容，请重新生成')
      // 不完整的输出不能进入正文替换流程。
      if (!completed) return ElMessage.warning('生成未完整结束，已保留结果供复制，请重新生成后再应用')

      resultReady.value = true
      history.value.push(
        { role: 'user', content: userContent.slice(0, 2000) },
        { role: 'assistant', content: generatedText.slice(0, 2000) }
      )
      if (history.value.length > 6) history.value = history.value.slice(-6)

      // 润色/改写且原文来自选区 → diff 直接铺在正文编辑器中
      if (snapshot) {
        const result = aiStore.openEditorDiff(snapshot, generatedText)
        if (!result.ok) ElMessage.warning(result.message)
      }
    })
  }

  const abortStream = () => {
    const handle = streamHandle
    streamHandle = null
    handle?.abort()
    resultReady.value = false
    streaming.value = false
  }

  // ---- 结果操作 ----
  const insertResult = () => {
    if (!resultReady.value || !ownsResult()) return
    if (!editorCtx.value?.insertAtCursor?.(resultText.value)) return ElMessage.warning('当前编辑器无法插入，请手动复制')
    ElMessage.success('已插入到光标处')
    reset()
  }

  const fillDescription = () => {
    if (!resultReady.value || !ownsResult()) return
    if (editorCtx.value?.fillDescription?.(resultText.value)) {
      ElMessage.success('已回填到文章摘要')
      reset()
    } else {
      ElMessage.warning('当前页面不支持回填，请手动复制')
    }
  }

  const applyTags = () => {
    if (!resultReady.value || !ownsResult() || !tagSuggestion.value) return
    if (editorCtx.value?.applySuggestion?.(tagSuggestion.value)) {
      ElMessage.success('已回填分类与标签')
      reset()
    } else {
      ElMessage.warning('当前页面不支持回填，请手动复制')
    }
  }

  const discardEditorDiff = () => {
    aiStore.closeEditorDiff()
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
    abortStream()
    resultOwner = null
    if (editorDiffOpen.value) {
      aiStore.closeEditorDiff()
    }
    resultText.value = ''
    currentAction.value = null
    tagSuggestion.value = null
  }
  // 文档身份变化同步取消任务，防止旧响应写入新文章的结果、历史或表单。
  watch(
    [() => editorCtx.value, () => editorCtx.value?.getEditorState?.()?.editorId,
      () => editorCtx.value?.getEditorState?.()?.documentId, () => editorCtx.value?.getEditorState?.()?.active],
    () => {
      reset()
      history.value = []
      instruction.value = ''
      syncSelection()
    },
    { flush: 'sync' }
  )
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

.editor-diff-notice {
  padding: 10px 12px;
  border: 1px solid #f3d19e;
  border-radius: 6px;
  background: #fdf6ec;
  color: #b88230;
  font-size: 13px;
  line-height: 1.7;

  .notice-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 8px;
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
