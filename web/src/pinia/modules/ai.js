import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'

/**
 * AI 平台全局状态：
 * - contexts: 各页面注册的可共享上下文（如编辑器句柄、表单回填回调）
 * - agent 注册表渲染 AiDock 时从这里判断能力是否可用
 */
export const useAiStore = defineStore('ai', () => {
  const contexts = reactive({})

  // AiDock 抽屉开关
  const dockVisible = ref(false)
  // 当前激活的 agent id
  const activeAgentId = ref('writing-assistant')

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
    registerContext,
    unregisterContext,
    hasContext,
    openDock,
    closeDock
  }
})
