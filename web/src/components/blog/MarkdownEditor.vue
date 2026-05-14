<template>
  <div class="markdown-editor" :class="{ 'is-fullscreen': fullscreen }">
    <div class="markdown-toolbar">
      <div class="toolbar-group">
        <el-tooltip
          v-for="tool in tools"
          :key="tool.key"
          :content="tool.label"
          placement="top"
        >
          <el-button
            class="tool-button"
            :icon="tool.icon"
            text
            @click="insertMarkdown(tool)"
          />
        </el-tooltip>
      </div>
      <div class="toolbar-group">
        <el-tooltip content="复制内容" placement="top">
          <el-button class="tool-button" :icon="CopyDocument" text @click="copyContent" />
        </el-tooltip>
        <el-tooltip :content="previewVisible ? '隐藏预览' : '显示预览'" placement="top">
          <el-button class="tool-button" :icon="View" text @click="previewVisible = !previewVisible" />
        </el-tooltip>
        <el-tooltip :content="fullscreen ? '退出全屏' : '全屏编辑'" placement="top">
          <el-button class="tool-button" :icon="FullScreen" text @click="fullscreen = !fullscreen" />
        </el-tooltip>
      </div>
    </div>

    <div class="markdown-body" :style="{ minHeight: editorHeight }">
      <div class="editor-pane" :class="{ 'is-alone': !previewVisible }">
        <textarea
          ref="textareaRef"
          v-model="value"
          class="markdown-textarea"
          :placeholder="placeholder"
          spellcheck="false"
          @keydown.tab.prevent="insertText('  ', '', '')"
        />
      </div>
      <div v-if="previewVisible" class="preview-pane">
        <div v-if="value" class="markdown-preview" v-html="previewHtml" />
        <div v-else class="preview-empty">Markdown 预览</div>
      </div>
    </div>

    <div class="markdown-footer">
      <span>{{ stats.characters }} 字符</span>
      <span>{{ stats.words }} 字</span>
      <span>{{ stats.lines }} 行</span>
    </div>
  </div>
</template>

<script setup>
  import { computed, nextTick, ref } from 'vue'
  import { Marked } from 'marked'
  import { markedHighlight } from 'marked-highlight'
  import hljs from 'highlight.js'
  import {
    ChatLineSquare,
    CopyDocument,
    EditPen,
    Files,
    FullScreen,
    Link,
    List,
    MagicStick,
    Picture,
    Tickets,
    View
  } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'

  const props = defineProps({
    modelValue: {
      type: String,
      default: ''
    },
    placeholder: {
      type: String,
      default: '请输入 Markdown 内容'
    },
    height: {
      type: [Number, String],
      default: 520
    }
  })

  const emit = defineEmits(['update:modelValue'])

  const textareaRef = ref()
  const previewVisible = ref(true)
  const fullscreen = ref(false)

  const marked = new Marked(
    {
      gfm: true,
      breaks: true
    },
    markedHighlight({
      langPrefix: 'hljs language-',
      highlight(code, lang) {
        const language = hljs.getLanguage(lang) ? lang : 'plaintext'
        return hljs.highlight(code, { language }).value
      }
    })
  )

  const value = computed({
    get: () => props.modelValue || '',
    set: (val) => emit('update:modelValue', val)
  })

  const editorHeight = computed(() => {
    if (typeof props.height === 'number') {
      return `${props.height}px`
    }
    return props.height
  })

  const previewHtml = computed(() => marked.parse(value.value || ''))

  const stats = computed(() => {
    const content = value.value || ''
    const words = content
      .replace(/```[\s\S]*?```/g, ' ')
      .replace(/[#>*_`~\-[\]()!|]/g, ' ')
      .match(/[\u4e00-\u9fa5]|[a-zA-Z0-9]+/g)

    return {
      characters: content.length,
      words: words ? words.length : 0,
      lines: content ? content.split(/\r?\n/).length : 0
    }
  })

  const tools = [
    { key: 'bold', label: '加粗', icon: EditPen, prefix: '**', suffix: '**', sample: '加粗文字' },
    { key: 'quote', label: '引用', icon: ChatLineSquare, block: '> 引用内容' },
    { key: 'list', label: '无序列表', icon: List, block: '- 列表项' },
    { key: 'code', label: '代码块', icon: Tickets, block: '```js\nconsole.log("hello")\n```' },
    { key: 'link', label: '链接', icon: Link, prefix: '[', suffix: '](https://)', sample: '链接文字' },
    { key: 'image', label: '图片', icon: Picture, prefix: '![', suffix: '](https://)', sample: '图片描述' },
    { key: 'table', label: '表格', icon: Files, block: '| 标题 | 内容 |\n| --- | --- |\n| 示例 | 文本 |' },
    { key: 'divider', label: '分割线', icon: MagicStick, block: '---' },
    { key: 'heading', label: '标题', icon: EditPen, block: '## 标题' }
  ]

  const focusTextarea = () => nextTick(() => textareaRef.value?.focus())

  const insertText = (prefix, suffix = '', sample = '') => {
    const textarea = textareaRef.value
    if (!textarea) return

    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const selected = value.value.slice(start, end)
    const text = `${prefix}${selected || sample}${suffix}`
    value.value = `${value.value.slice(0, start)}${text}${value.value.slice(end)}`

    nextTick(() => {
      const cursorStart = start + prefix.length
      const cursorEnd = cursorStart + (selected || sample).length
      textarea.setSelectionRange(cursorStart, cursorEnd)
      textarea.focus()
    })
  }

  const insertBlock = (block) => {
    const textarea = textareaRef.value
    if (!textarea) return

    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const before = value.value.slice(0, start)
    const after = value.value.slice(end)
    const selected = value.value.slice(start, end)
    const content = selected || block
    const leading = before && !before.endsWith('\n') ? '\n' : ''
    const trailing = after && !after.startsWith('\n') ? '\n' : ''
    const text = `${leading}${content}\n${trailing}`

    value.value = `${before}${text}${after}`

    nextTick(() => {
      const cursorStart = start + leading.length
      const cursorEnd = cursorStart + content.length
      textarea.setSelectionRange(cursorStart, cursorEnd)
      textarea.focus()
    })
  }

  const insertMarkdown = (tool) => {
    if (tool.block) {
      insertBlock(tool.block)
      return
    }
    insertText(tool.prefix, tool.suffix, tool.sample)
  }

  const copyContent = async () => {
    try {
      await navigator.clipboard.writeText(value.value)
      ElMessage.success('已复制')
    } catch (error) {
      ElMessage.error('复制失败')
    }
  }

  defineExpose({ focusTextarea })
</script>

<style scoped lang="scss">
.markdown-editor {
  width: 100%;
  box-sizing: border-box;
  overflow: hidden;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  background: #fff;
}

.markdown-editor.is-fullscreen {
  position: fixed;
  z-index: 3000;
  inset: 16px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 12px 36px rgb(0 0 0 / 16%);

  .markdown-body {
    flex: 1;
  }
}

.markdown-toolbar,
.markdown-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  background: #f7f8fa;
}

.markdown-toolbar {
  border-bottom: 1px solid #ebeef5;
}

.markdown-footer {
  gap: 16px;
  justify-content: flex-end;
  border-top: 1px solid #ebeef5;
  color: #909399;
  font-size: 12px;
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 2px;
}

.tool-button {
  width: 30px;
  height: 30px;
  padding: 0;
}

.markdown-body {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  background: #fff;
}

.editor-pane,
.preview-pane {
  min-width: 0;
}

.editor-pane {
  border-right: 1px solid #ebeef5;
}

.editor-pane.is-alone {
  grid-column: 1 / -1;
  border-right: 0;
}

.markdown-textarea {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  min-height: inherit;
  padding: 18px;
  border: 0;
  outline: none;
  resize: none;
  color: #303133;
  background: #fff;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace;
  font-size: 14px;
  line-height: 1.75;
}

.preview-pane {
  overflow: auto;
  background: #fff;
}

.preview-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #c0c4cc;
}

.markdown-preview {
  padding: 18px 22px;
  color: #303133;
  line-height: 1.8;

  :deep(h1),
  :deep(h2),
  :deep(h3) {
    margin: 20px 0 12px;
    color: #1f2d3d;
    line-height: 1.35;
  }

  :deep(h1) {
    font-size: 26px;
  }

  :deep(h2) {
    padding-bottom: 6px;
    border-bottom: 1px solid #ebeef5;
    font-size: 22px;
  }

  :deep(h3) {
    font-size: 18px;
  }

  :deep(p),
  :deep(ul),
  :deep(ol),
  :deep(blockquote),
  :deep(pre),
  :deep(table) {
    margin: 0 0 14px;
  }

  :deep(ul),
  :deep(ol) {
    padding-left: 1.6em;
  }

  :deep(ul) {
    list-style: disc;
  }

  :deep(ol) {
    list-style: decimal;
  }

  :deep(li) {
    margin: 4px 0;
  }

  :deep(li > ul) {
    margin: 4px 0 0;
    list-style: circle;
  }

  :deep(li > ol) {
    margin: 4px 0 0;
  }

  :deep(.contains-task-list) {
    padding-left: 1.2em;
    list-style: none;
  }

  :deep(input[type='checkbox']) {
    margin-right: 6px;
  }

  :deep(blockquote) {
    padding: 8px 14px;
    border-left: 4px solid #409eff;
    background: #f5f7fa;
    color: #606266;
  }

  :deep(pre) {
    overflow: auto;
    padding: 14px;
    border-radius: 6px;
    background: #f6f8fa;
  }

  :deep(code) {
    padding: 2px 5px;
    border-radius: 4px;
    background: #f6f8fa;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace;
    font-size: 13px;
  }

  :deep(pre code) {
    padding: 0;
    background: transparent;
  }

  :deep(.hljs) {
    color: #24292e;
    background: #f6f8fa;
  }

  :deep(.hljs-doctag),
  :deep(.hljs-keyword),
  :deep(.hljs-meta .hljs-keyword),
  :deep(.hljs-template-tag),
  :deep(.hljs-template-variable),
  :deep(.hljs-type),
  :deep(.hljs-variable.language_) {
    color: #d73a49;
  }

  :deep(.hljs-title),
  :deep(.hljs-title.class_),
  :deep(.hljs-title.class_.inherited__),
  :deep(.hljs-title.function_) {
    color: #6f42c1;
  }

  :deep(.hljs-attr),
  :deep(.hljs-attribute),
  :deep(.hljs-literal),
  :deep(.hljs-meta),
  :deep(.hljs-number),
  :deep(.hljs-operator),
  :deep(.hljs-selector-attr),
  :deep(.hljs-selector-class),
  :deep(.hljs-selector-id),
  :deep(.hljs-variable) {
    color: #005cc5;
  }

  :deep(.hljs-regexp),
  :deep(.hljs-string),
  :deep(.hljs-meta .hljs-string) {
    color: #032f62;
  }

  :deep(.hljs-built_in),
  :deep(.hljs-symbol) {
    color: #e36209;
  }

  :deep(.hljs-code),
  :deep(.hljs-comment),
  :deep(.hljs-formula) {
    color: #6a737d;
  }

  :deep(img) {
    max-width: 100%;
    border-radius: 6px;
  }

  :deep(table) {
    width: 100%;
    border-collapse: collapse;
  }

  :deep(th),
  :deep(td) {
    padding: 8px 10px;
    border: 1px solid #dcdfe6;
  }
}

@media (max-width: 960px) {
  .markdown-body {
    grid-template-columns: 1fr;
  }

  .editor-pane {
    border-right: 0;
    border-bottom: 1px solid #ebeef5;
  }

  .preview-pane {
    min-height: 320px;
  }
}
</style>
