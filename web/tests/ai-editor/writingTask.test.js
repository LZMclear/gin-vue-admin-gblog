import { test } from 'node:test'
import assert from 'node:assert/strict'
import { parseTitleCandidates, createWritingTask, retryTaskError, titleFillError } from '../../src/components/ai/agents/writing-assistant/writingTask.js'
import { captureEditorSnapshot } from '../../src/components/blog/editorSnapshot.js'

test('标题候选清理常见列表格式、去重并限制为五项', () => {
  assert.deepEqual(parseTitleCandidates('候选标题：\n1. **标题一**\n2、标题二\n- “标题三”\n标题一\n<script>\n' + '长'.repeat(121)), ['标题一', '标题二', '标题三'])
  assert.equal(parseTitleCandidates('一\n二\n三\n四\n五\n六').length, 5)
})

test('回填标题不覆盖生成之后人工改过的标题', () => {
  assert.equal(titleFillError('新标题', '旧标题', '旧标题'), '')
  assert.match(titleFillError('新标题', '旧标题', '人工标题'), /已被编辑/)
  for (const title of ['', '长'.repeat(121), '两\n行', '<tag>']) assert.notEqual(titleFillError(title, '', ''), '')
})

test('重试固定请求条件和历史，不采用后来变动的输入', () => {
  const state = { active: true, editorId: 'one', documentId: 'a', revision: 1, content: '原正文' }
  let title = '旧标题'
  const context = { getEditorState: () => state, getFullText: () => state.content, getTitle: () => title }
  const history = [{ role: 'user', content: '原历史' }]
  const payload = { content: state.content, title, tone: 'formal', instruction: '原指令' }
  const task = createWritingTask(payload, null, context, history)
  history[0].content = '后改历史'; payload.instruction = '后改指令'
  assert.equal(task.payload.history[0].content, '原历史')
  assert.equal(task.payload.instruction, '原指令')
  assert.equal(retryTaskError(task, context), '')
  assert.match(retryTaskError(task, { ...context }), /切换/)
  state.content = '改后正文'
  assert.match(retryTaskError(task, context), /已变化/)
  state.content = '原正文'; title = '新标题'
  assert.match(retryTaskError(task, context), /已变化/)
})

test('选区重试及自定义对比保留原始快照，不因当前选区改变重定位', () => {
  const state = { active: true, editorId: 'one', documentId: 'a', revision: 1, content: '前文原选区后文' }
  const context = { getEditorState: () => state, getFullText: () => state.content, getTitle: () => '' }
  const snapshot = captureEditorSnapshot(state, { start: 2, end: 5 })
  const task = createWritingTask({ content: state.content, title: '', selection: snapshot.text }, snapshot, context)
  assert.equal(retryTaskError(task, context), '')
  assert.equal(task.snapshot.start, 2)
  state.revision++
  assert.match(retryTaskError(task, context), /正文已发生变化/)
})
