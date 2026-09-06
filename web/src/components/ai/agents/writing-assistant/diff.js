import { diffArrays } from 'diff'

/**
 * 将 Markdown 文本切分为块：代码围栏视为原子块，其余按空行分段。
 */
export function splitMarkdownBlocks(text) {
  const blocks = []
  const lines = String(text || '').split(/\r?\n/)
  let current = []
  let inFence = false

  const flush = () => {
    if (current.length) {
      blocks.push(current.join('\n'))
      current = []
    }
  }

  for (const line of lines) {
    if (/^\s*(```|~~~)/.test(line)) {
      // 围栏开始/结束行
      if (inFence) {
        current.push(line)
        flush() // 整个围栏作为一块
        inFence = false
      } else {
        flush()
        current.push(line)
        inFence = true
      }
      continue
    }
    if (!inFence && line.trim() === '') {
      flush()
      continue
    }
    current.push(line)
  }
  flush()
  return blocks
}

/**
 * 段落级 diff：返回块序列 [{ type: 'equal'|'modified'|'added'|'removed', original, revised, takeRevised }]
 * - equal: 两边一致
 * - modified: 两侧对应但内容不同（可逐块选择采用 AI 版本或保留原文）
 * - added: AI 新增块；removed: 仅原文有块
 */
export function diffMarkdownBlocks(originalText, revisedText) {
  const original = splitMarkdownBlocks(originalText)
  const revised = splitMarkdownBlocks(revisedText)

  const parts = diffArrays(original, revised, {
    comparator: (left, right) => left === right
  })

  const result = []
  for (const part of parts) {
    const count = part.count || (part.value || []).length
    if (part.added) {
      for (const block of part.value) {
        result.push({ type: 'added', original: null, revised: block, takeRevised: true })
      }
      void count
    } else if (part.removed) {
      for (const block of part.value) {
        result.push({ type: 'removed', original: block, revised: null, takeRevised: false })
      }
    } else {
      // 未变化
      for (const block of part.value) {
        result.push({ type: 'equal', original: block, revised: block, takeRevised: true })
      }
    }
  }

  // 把相邻的 added/removed 对合并为 modified（jsdiff 对修改会输出 -removed +added 相邻对）
  const merged = []
  for (let i = 0; i < result.length; i++) {
    const cur = result[i]
    const next = result[i + 1]
    if (cur.type === 'removed' && next && next.type === 'added') {
      merged.push({
        type: 'modified',
        original: cur.original,
        revised: next.revised,
        takeRevised: true
      })
      i++
      continue
    }
    merged.push(cur)
  }
  return merged
}

/**
 * 按 diff 块的用户选择拼回最终文本。
 */
export function applyDiffBlocks(blocks) {
  return blocks
    .map((block) => {
      if (block.type === 'equal') return block.original
      if (block.type === 'modified' || block.type === 'added') {
        return block.takeRevised ? block.revised : block.original
      }
      // removed：takeRevised=false 表示保留原文
      return block.takeRevised ? null : block.original
    })
    .filter((text) => text !== null && text !== undefined && text !== '')
    .join('\n\n')
}

/**
 * 统计可选择的块数量与已采纳数量。
 */
export function diffStats(blocks) {
  const selectable = blocks.filter(
    (b) => b.type === 'modified' || b.type === 'added' || b.type === 'removed'
  )
  const adopted = selectable.filter((b) => b.takeRevised)
  return { total: selectable.length, adopted: adopted.length }
}
