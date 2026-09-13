import { test } from 'node:test'
import assert from 'node:assert/strict'
import { buildSuggestionPatch, cursorContext } from '../../src/components/ai/agents/writing-assistant/suggestion.js'

const categories = [{ id: 1, categoryName: '技术' }]
const tags = [{ id: 2, tagName: 'Vue' }, { id: 3, tagName: 'Go' }]
test('推荐回填保留已有标签，按 ID 追加并去重', () => {
  const form = { cate: 9, tagList: [3] }
  const result = buildSuggestionPatch(form, { category: '技术', categoryId: 1, tags: ['Vue', 'Go'], tagIds: [2, 3], newTags: [] }, categories, tags)
  assert.deepEqual(result, { ok: true, patch: { cate: 1, tagList: [3, 2] } })
  assert.deepEqual(form, { cate: 9, tagList: [3] })
})
test('不存在的已有项和过期 ID 不得变成新建字符串，变更原子化', () => {
  for (const suggestion of [{ category: '未知' }, { tags: ['未知'] }, { tags: ['Vue'], tagIds: [99] }]) {
    assert.equal(buildSuggestionPatch({ tagList: [3] }, suggestion, categories, tags).ok, false)
  }
})
test('只有明确选择的新标签才作为字符串回填；已存在的新建议使用 ID', () => {
  assert.deepEqual(buildSuggestionPatch({ tagList: [3] }, { newTags: ['新标签', ' 新标签 ', 'vue'] }, categories, tags).patch.tagList, [3, '新标签', 2])
  assert.equal(buildSuggestionPatch({}, { newTags: ['<bad>'] }, categories, tags).ok, false)
  assert.equal(buildSuggestionPatch({}, { tags: [], newTags: [] }, categories, tags).ok, false)
})
test('光标在零位置不会读取全文，emoji 使用浏览器 UTF-16 偏移', () => {
  assert.equal(cursorContext('甲😀乙', 0), '')
  assert.equal(cursorContext('甲😀乙', 3), '甲😀')
  assert.equal(cursorContext('甲😀乙'), '甲😀乙')
})
