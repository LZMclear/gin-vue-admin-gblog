import assert from 'node:assert/strict'
import { test } from 'node:test'
import { captureEditorSnapshot, applySnapshot, undoSnapshot, textareaOffsetToSource, sourceOffsetToTextarea } from '../../src/components/blog/editorSnapshot.js'
import { createPinia, setActivePinia } from 'pinia'
import { useAiStore } from '../../src/pinia/modules/ai.js'

const original = '前文\r\n\r\n待修改😀\r\n\r\n后文\r\n'
const state = { active: true, editorId: 'editor-a', documentId: 'article-a', revision: 3, content: original }
const start = original.indexOf('待修改')
const end = original.indexOf('\r\n', start)

for (const [from, to] of [[0, 2], [start, end], [original.indexOf('后文'), original.length], [0, original.length]]) {
  test(`固定选区替换及撤销 ${from}:${to}`, () => {
    const snapshot = captureEditorSnapshot(state, { start: from, end: to })
    const result = applySnapshot(state, snapshot, 'AI')
    assert.equal(result.ok, true)
    assert.equal(result.content, original.slice(0, from) + 'AI' + original.slice(to))
    const restored = undoSnapshot({ ...state, content: result.content }, result.undo)
    assert.equal(restored.content, original)
  })
}

test('无选区或无编辑器时不可将状态误解为整篇替换', () => {
  assert.equal(captureEditorSnapshot(state, { start: 0, end: 0 }), null)
  assert.equal(captureEditorSnapshot({ ...state, active: false }, { start, end }), null)
  assert.equal(applySnapshot(state, null, 'AI').ok, false)
})

test('正文、文章、实例或版本改变后拒绝应用，离开文章后也拒绝', () => {
  const snapshot = captureEditorSnapshot(state, { start, end })
  for (const change of [{ content: original + '人工编辑' }, { revision: 4 }, { documentId: 'article-b' }, { editorId: 'editor-b' }, { active: false }]) {
    assert.equal(applySnapshot({ ...state, ...change }, snapshot, 'AI').ok, false)
  }
})

test('撤销不覆盖后来输入的内容，不跨文章恢复', () => {
  const snapshot = captureEditorSnapshot(state, { start, end })
  const result = applySnapshot(state, snapshot, 'AI')
  assert.equal(undoSnapshot({ ...state, content: result.content + '新输入' }, result.undo).ok, false)
  assert.equal(undoSnapshot({ ...state, content: result.content, documentId: 'article-b' }, result.undo).ok, false)
})

test('CRLF、中文和 emoji 的 DOM 偏移正确映射回原文', () => {
  const displayed = original.replace(/\r\n?/g, '\n')
  const domStart = displayed.indexOf('待修改')
  const domEnd = displayed.indexOf('\n', domStart)
  assert.equal(textareaOffsetToSource(original, domStart), start)
  assert.equal(textareaOffsetToSource(original, domEnd), end)
  assert.equal(sourceOffsetToTextarea(original, end), domEnd)
})

test('Store 预览保留选区前后文，固定范围应用失败时保留对比', () => {
  setActivePinia(createPinia())
  const store = useAiStore()
  store.registerContext('editor', { getEditorState: () => state })
  const snapshot = captureEditorSnapshot(state, { start, end })
  assert.equal(store.openEditorDiff(snapshot, 'AI').ok, true)
  assert.equal(store.diffPreviewText, original.slice(0, start) + 'AI' + original.slice(end))
  assert.equal(store.applyEditorDiff({}).ok, false)
  assert.equal(store.diff.active, true)
  assert.equal(store.applyEditorDiff({ applySelectionSnapshot: (saved, text) => applySnapshot(state, saved, text) }).ok, true)
  assert.equal(store.diff.active, false)
})

test('旧页面注销不会清掉新编辑器，迟到的结果不能绑定新文章', () => {
  setActivePinia(createPinia())
  const store = useAiStore()
  const oldHandle = { getEditorState: () => state }
  const newHandle = { getEditorState: () => ({ ...state, editorId: 'editor-b' }) }
  store.registerContext('editor', oldHandle)
  store.registerContext('editor', newHandle)
  store.unregisterContext('editor', oldHandle)
  assert.equal(store.contexts.editor, newHandle)
  assert.equal(store.openEditorDiff(captureEditorSnapshot(state, { start, end }), 'AI').ok, false)
})
