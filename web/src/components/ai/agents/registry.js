import WritingAssistantPanel from './writing-assistant/WritingAssistantPanel.vue'

/**
 * AI Agent 注册表：新增 agent 时在此追加一项即可，AiDock 自动渲染。
 * - needContext: agent 依赖的页面上下文（对应 aiStore.registerContext 注册的名称）
 */
export const agentRegistry = [
  {
    id: 'writing-assistant',
    name: '写作助手',
    icon: 'EditPen',
    description: '润色、改写、续写、大纲、摘要、标签推荐',
    component: WritingAssistantPanel,
    needContext: ['editor']
  }
]
