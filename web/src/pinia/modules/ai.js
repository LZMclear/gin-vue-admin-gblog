import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'
import { diffMarkdownBlocks, applyDiffBlocks } from '@/components/ai/agents/writing-assistant/diff'

/**
 * AI 平台全局状态：
 * - contexts: 各页面注册的可共享上下文（如编辑器句柄、表单回填回调）
 * - diff: 润色/改写的编辑器内 diff 状态（MarkdownEditor 开启 enable-ai-diff 后展示）
 * - agent 注册表渲染 AiDock 时从这里判断能力是否可用
 */
export const useAiStore = defineStore('ai', () => {
  const contexts = reactive({})

  // AiDock 抽屉开关
  const dockVisible = ref(false)
  // 当前激活的 agent id
  const activeAgentId = ref('writing-assistant')

  // 编辑器内 diff 状态
  const diff = reactive({ active: false, blocks: [], sourceText: '' })

  const mergedDiffText = computed(() =>
    diff.active ? applyDiffBlocks(diff.blocks) : ''
  )

  const openEditorDiff = (originalText, resultText) => {
    diff.blocks = diffMarkdownBlocks(originalText, resultText)
    diff.sourceText = originalText || ''
    diff.active = true
  }

  const closeEditorDiff = () => {
    diff.active = false
    diff.blocks = []
    diff.sourceText = ''
  }

  // 把逐块采纳后的最终文本写回编辑器选区；handle 为编辑器句柄
  const applyEditorDiff = (handle) => {
    if (!diff.active) return false
    const finalText = applyDiffBlocks(diff.blocks)
    handle?.replaceSelection?.(finalText)
    closeEditorDiff()
    return true
  }

  const registerContext = (name, handle) => {
    contexts[name] = handle
  }

  const unregisterContext = (name) => {
    delete contexts[name]
  }

  const hasContext = (name) => computed(() => Boolean(contexts[name]))

  const openDock = (agentId) => {
    if (agentId) activeAgentId.value = agentId
    dockVisible.value = true
  }

  const closeDock = () => {
    dockVisible.value = false
  }

  return {
    contexts,
    dockVisible,
    activeAgentId,
    diff,
    mergedDiffText,
    openEditorDiff,
    closeEditorDiff,
    applyEditorDiff,
    registerContext,
    unregisterContext,
    hasContext,
    openDock,
    closeDock
  }
})
