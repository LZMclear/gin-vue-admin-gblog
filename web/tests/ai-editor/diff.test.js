import assert from 'node:assert/strict'
import { test } from 'node:test'
import { applyDiffBlocks, diffMarkdownBlocks, splitMarkdownBlocks } from '../../src/components/ai/agents/writing-assistant/diff.js'

test('两段连续改写按顺序对应，默认采用不会保留多余旧段落', () => {
  const blocks = diffMarkdownBlocks('旧 A\n\n旧 B', '新 A\n\n新 B')
  assert.deepEqual(blocks.map((b) => [b.type, b.original, b.revised]), [
    ['modified', '旧 A\n\n', '新 A\n\n'], ['modified', '旧 B', '新 B']
  ])
  assert.equal(applyDiffBlocks(blocks), '新 A\n\n新 B')
  blocks[0].takeRevised = false
  assert.equal(applyDiffBlocks(blocks), '旧 A\n\n新 B')
})

const cases = [
  ['', '新内容'], ['旧内容', ''], ['', ''],
  ['\r\n A\r\n\r\n\r\n B\r\n', '\n 新 A\n\n 新 B\n'],
  ['A\n\nB\n\nC', 'A\n\n新 B\n\n新增\n\nC'],
  ['A\n\n删除\n\nB', 'A\n\nB'],
  ['A', 'A\n\nB'], ['A\n\nB', 'B\n\nA'],
  ['\t\r\n\r\n', '\n'],
  ['[文档][ref]\n\n[ref]: /docs\n', '[资料][ref]\n\n[ref]: /new\n']
]
for (const [original, revised] of cases) {
  test(`两侧无损还原 ${JSON.stringify(original)} → ${JSON.stringify(revised)}`, () => {
    const blocks = diffMarkdownBlocks(original, revised)
    assert.equal(applyDiffBlocks(blocks), revised)
    blocks.forEach((b) => { b.takeRevised = false })
    assert.equal(applyDiffBlocks(blocks), original)
    assert.equal(applyDiffBlocks(diffMarkdownBlocks(original, original)), original)
  })
}

test('段落数量改变时整组采纳，不猜测段落对应关系', () => {
  const blocks = diffMarkdownBlocks('第一段\n\n第二段', '合并后的一段')
  assert.equal(blocks.length, 1)
  assert.equal(blocks[0].fallback, true)
})

test('代码围栏、列表、表格和引用保持完整', () => {
  const blocks = [
    '````markdown\r\n```js\r\nconst x = 1\r\n```\r\n~~~~\r\n````\r\n\r\n',
    '- 一\r\n\r\n  续段\r\n- 二\r\n\r\n',
    '| A | B |\r\n| --- | --- |\r\n| 1 | 2 |\r\n\r\n',
    '> 引用\r\n>\r\n> 引用的第二段\r\n'
  ]
  assert.deepEqual(splitMarkdownBlocks(blocks.join('')), blocks)
})

test('引用定义与未闭合围栏也不丢失原文', () => {
  for (const text of ['[ref]: /docs\r\n\r\n正文', '~~~js\n内容\n\n仍在代码块', '    缩进代码\n\n    第二行\n']) {
    assert.equal(splitMarkdownBlocks(text).join(''), text)
  }
})
