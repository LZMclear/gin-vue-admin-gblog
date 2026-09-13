<template>
  <ModelConfig v-if="showModelConfig" />
  <div v-else class="fixture">
    <MarkdownEditor ref="editor" v-model="content" :document-id="documentId" enable-ai-diff />
    <WritingAssistantPanel />
  </div>
</template>

<script setup>
  import { nextTick, onMounted, ref } from 'vue'
  import MarkdownEditor from '../../src/components/blog/MarkdownEditor.vue'
  import WritingAssistantPanel from '../../src/components/ai/agents/writing-assistant/WritingAssistantPanel.vue'
  import ModelConfig from '../../src/view/ai/modelConfig/modelConfig.vue'
  import { useAiStore } from '../../src/pinia/modules/ai.js'
  import { renderSafeMarkdown } from '../../src/utils/safeMarkdown.js'

  const editor = ref()
  const showModelConfig = new URLSearchParams(window.location.search).has('models')
  const content = ref('前文\n\n选中内容\n\n后文')
  const documentId = ref('article-a')
  const description = ref('原摘要')
  const suggestion = ref(null)
  const store = useAiStore()
  onMounted(() => {
    if (showModelConfig) return
    store.registerContext('editor', {
      getSelection: () => editor.value.getSelection(),
      getEditorState: () => editor.value.getEditorState(),
      captureSelection: () => editor.value.captureSelection(),
      applySelectionSnapshot: (snapshot, text) => editor.value.applySelectionSnapshot(snapshot, text),
      getFullText: () => content.value,
      insertAtCursor: (text) => editor.value.insertAtCursor(text),
      fillDescription: (text) => { description.value = text; return true },
      applySuggestion: (value) => { suggestion.value = value; return true },
      getTitle: () => '回归测试文章'
    })
    window.aiEditorTest = {
      getContent: () => content.value,
      setContent: async (text) => { content.value = text; await nextTick() },
      changeDocument: async () => { documentId.value = 'article-b'; await nextTick() },
      getDescription: () => description.value,
      getSuggestion: () => suggestion.value,
      leaveEditor: () => store.unregisterContext('editor'),
      renderSafeMarkdown
    }
  })
</script>

<style>
  body { margin: 20px; }
  .fixture { display: grid; grid-template-columns: minmax(0, 1fr) 380px; gap: 20px; }
</style>
