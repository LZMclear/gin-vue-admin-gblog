import assert from 'node:assert/strict'
import { test } from 'node:test'
import { createPinia, setActivePinia } from 'pinia'
import { useAiStore } from '../../src/pinia/modules/ai.js'
import { captureEditorSnapshot } from '../../src/components/blog/editorSnapshot.js'
import { placeSelectionToolbar } from '../../src/components/blog/selectionToolbarPosition.js'

const state = { active: true, editorId: 'editor-a', documentId: 'article-a', revision: 1, content: '前文\n选中内容😀\n后文' }
const snapshot = captureEditorSnapshot(state, { start: 3, end: 9 })
function setup() {
  setActivePinia(createPinia())
  const store = useAiStore()
  store.registerContext('editor', { getEditorState: () => state })
  return store
}

test('选区入口打开写作助手并保留点击时的快照，重复点击不覆盖请求', () => {
  const store = setup()
  assert.equal(store.requestSelectionAction('polish', snapshot).ok, true)
  assert.equal(store.dockVisible, true)
  assert.equal(store.activeAgentId, 'writing-assistant')
  assert.equal(store.selectionAction.snapshot, snapshot)
  assert.equal(store.requestSelectionAction('rewrite', snapshot).ok, false)
  assert.equal(store.selectionAction.action, 'polish')
})

test('生成中、对比中、空选区和失效选区不发起新任务', () => {
  const store = setup()
  assert.equal(store.requestSelectionAction('polish', null).ok, false)
  assert.equal(store.requestSelectionAction('polish', { ...snapshot, revision: 0 }).ok, false)
  store.writingBusy = true
  assert.equal(store.requestSelectionAction('rewrite', snapshot).ok, false)
  store.writingBusy = false
  store.openEditorDiff(snapshot, '修改')
  assert.equal(store.requestSelectionAction('rewrite', snapshot).ok, false)
  assert.equal(store.selectionAction, null)
  assert.equal(store.dockVisible, false)
})

test('编辑器切换或注销清理待处理选区，旧页面注销不清理新任务', () => {
  const store = setup()
  const oldHandle = store.contexts.editor
  store.requestSelectionAction('polish', snapshot)
  const newState = { ...state, editorId: 'editor-b', documentId: 'article-b' }
  const newHandle = { getEditorState: () => newState }
  store.registerContext('editor', newHandle)
  assert.equal(store.selectionAction, null)
  const newSnapshot = captureEditorSnapshot(newState, { start: 3, end: 9 })
  store.requestSelectionAction('rewrite', newSnapshot)
  store.unregisterContext('editor', oldHandle)
  assert.equal(store.selectionAction.snapshot, newSnapshot)
  store.unregisterContext('editor', newHandle)
  assert.equal(store.selectionAction, null)
})

test('选区入口在左右边界保持可见，顶部空间不足时移到下方', () => {
  const viewport = { width: 800, height: 600 }
  assert.deepEqual(placeSelectionToolbar({ x: 400, y: 100, height: 20 }, viewport), { left: 308, top: 54 })
  assert.equal(placeSelectionToolbar({ x: 0, y: 100, height: 20 }, viewport).left, 8)
  assert.equal(placeSelectionToolbar({ x: 800, y: 100, height: 20 }, viewport).left, 608)
  assert.equal(placeSelectionToolbar({ x: 400, y: 10, height: 20 }, viewport).top, 38)
})
