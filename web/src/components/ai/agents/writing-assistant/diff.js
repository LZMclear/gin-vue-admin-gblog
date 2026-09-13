import { diffArrays } from 'diff'
import { Lexer } from 'marked'

// 列表、表格、引用及完整代码围栏作为原子块，raw 映射回原文以保留 CRLF 和空白。
function tokenizeMarkdown(text) {
  const source = String(text || '')
  if (!source) return []
  const normalized = source.replace(/\r\n?/g, '\n')
  const tokens = Lexer.lex(normalized, { gfm: true })
  // 引用定义等语法可能不出现在顶层 token 中，整段比较以保证无损。
  if (tokens.map((token) => token.raw).join('') !== normalized) {
    return [{ kind: 'document', raw: source }]
  }
  const blocks = []
  let offset = 0
  let leading = ''
  for (const token of tokens) {
    const start = offset
    for (let i = 0; i < token.raw.length; i++) {
      offset += source[offset] === '\r' && source[offset + 1] === '\n' ? 2 : 1
    }
    const raw = source.slice(start, offset)
    if (token.type === 'space') {
      if (blocks.length) blocks[blocks.length - 1].raw += raw
      else leading += raw
    } else {
      blocks.push({ kind: token.type, raw: leading + raw })
      leading = ''
    }
  }
  if (leading) blocks.push({ kind: 'space', raw: leading })
  return blocks
}

export function splitMarkdownBlocks(text) {
  return tokenizeMarkdown(text).map((block) => block.raw)
}

function changeBlock(original, revised, fallback = false) {
  return {
    type: original === revised ? 'equal' : original === null ? 'added' : revised === null ? 'removed' : 'modified',
    original,
    revised,
    takeRevised: true,
    fallback
  }
}

/** 相同块作锚点；连续修改成组处理，无法逐段对应时作为一组采纳。 */
export function diffMarkdownBlocks(originalText, revisedText) {
  const original = tokenizeMarkdown(originalText)
  const revised = tokenizeMarkdown(revisedText)
  const parts = diffArrays(original, revised, {
    comparator: (left, right) => left.kind === right.kind && left.raw === right.raw,
    timeout: 100
  })
  if (!parts) {
    return [changeBlock(String(originalText || ''), String(revisedText || ''), true)]
  }
  const result = []
  for (let i = 0; i < parts.length;) {
    const part = parts[i]
    if (!part.added && !part.removed) {
      result.push(...part.value.map((block) => changeBlock(block.raw, block.raw)))
      i++
      continue
    }
    const removed = []
    const added = []
    while (i < parts.length && (parts[i].added || parts[i].removed)) {
      const current = parts[i++]
      const group = current.removed ? removed : added
      group.push(...current.value)
    }
    const canAlign = removed.length === added.length && removed.every(
      (block, index) => block.kind === added[index].kind && block.kind !== 'document'
    )
    if (canAlign) {
      result.push(...removed.map((block, index) => changeBlock(block.raw, added[index].raw)))
    } else {
      // 增删段、结构改变时不猜测对应关系，避免局部采纳破坏 Markdown。
      result.push(changeBlock(
        removed.length ? removed.map((block) => block.raw).join('') : null,
        added.length ? added.map((block) => block.raw).join('') : null,
        removed.length > 0 && added.length > 0
      ))
    }
  }
  return result
}

/** 拼接选定的原始片段，不插入、删除或归一化任何空白。 */
export function applyDiffBlocks(blocks) {
  return blocks.map((block) => (
    block.type === 'equal' || !block.takeRevised ? block.original : block.revised
  ) ?? '').join('')
}

export function diffStats(blocks) {
  const selectable = blocks.filter((block) => block.type !== 'equal')
  return { total: selectable.length, adopted: selectable.filter((block) => block.takeRevised).length }
}
