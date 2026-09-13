import { defineStore } from 'pinia'
import { computed, markRaw, reactive, ref } from 'vue'
import { diffMarkdownBlocks, applyDiffBlocks } from '../../components/ai/agents/writing-assistant/diff.js'
import { snapshotError, replaceSnapshotRange } from '../../components/blog/editorSnapshot.js'

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
  const diff = reactive({ active: false, blocks: [], snapshot: null })

  const mergedDiffText = computed(() =>
    diff.active ? applyDiffBlocks(diff.blocks) : ''
  )

  const diffPreviewText = computed(() => diff.active && diff.snapshot
    ? replaceSnapshotRange(diff.snapshot, mergedDiffText.value)
    : '')

  const openEditorDiff = (snapshot, resultText) => {
    const message = snapshotError(snapshot, contexts.editor?.getEditorState?.())
    if (message) return { ok: false, message }
    if (diff.active) return { ok: false, message: '请先应用或取消当前对比' }
    diff.blocks = diffMarkdownBlocks(snapshot.text, resultText)
    diff.snapshot = snapshot
    diff.active = true
    return { ok: true }
  }

  const closeEditorDiff = () => {
    diff.active = false
    diff.blocks = []
    diff.snapshot = null
  }

  const setDiffChoice = ({ index, takeRevised }) => {
    const blocks = Number.isInteger(index) ? [diff.blocks[index]] : diff.blocks
    for (const block of blocks) {
      if (block && block.type !== 'equal') block.takeRevised = takeRevised
    }
  }

  // 写回发起请求时的固定选区；失败时保留对比，不触碰正文。
  const applyEditorDiff = (handle) => {
    if (!diff.active) return { ok: false, message: '没有待应用的 AI 修改' }
    const finalText = applyDiffBlocks(diff.blocks)
    const result = handle?.applySelectionSnapshot?.(diff.snapshot, finalText)
    if (!result?.ok) return result || { ok: false, message: '当前编辑器无法应用修改' }
    closeEditorDiff()
    return result
  }

  const registerContext = (name, handle) => {
    contexts[name] = markRaw(handle)
  }

  const unregisterContext = (name, handle) => {
    if (handle && contexts[name] !== handle) return
    if (name === 'editor' && diff.snapshot?.editorId === contexts[name]?.getEditorState?.()?.editorId) {
      closeEditorDiff()
    }
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
    diffPreviewText,
    openEditorDiff,
    closeEditorDiff,
    setDiffChoice,
    applyEditorDiff,
    registerContext,
    unregisterContext,
    hasContext,
    openDock,
    closeDock
  }
})
