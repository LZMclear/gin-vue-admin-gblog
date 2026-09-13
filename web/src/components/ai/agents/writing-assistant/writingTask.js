import { snapshotError } from '../../../blog/editorSnapshot.js'

export function parseTitleCandidates(text) {
  const titles = String(text || '').split(/\r?\n/).map(line => line.trim()
    .replace(/^(?:#{1,6}\s+|[-*+]\s+|\d+[.)、]\s*|[一二三四五六七八九十]+[、.)]\s*)/, '')
    .replace(/^\*\*(.+)\*\*$/, '$1').replace(/^[“"「](.*)[”"」]$/, '$1').trim())
    .filter(line => line && [...line].length <= 120 && !/[<>`]/.test(line) && !/[：:]$/.test(line))
  return [...new Set(titles)].slice(0, 5)
}

export function createWritingTask(payload, snapshot, context, history = [], autoDiff = false) {
  const state = context?.getEditorState?.()
  return {
    context, editorId: state?.editorId, documentId: state?.documentId,
    payload: { ...payload, history: history.map(message => ({ ...message })) }, snapshot, autoDiff
  }
}

export function retryTaskError(task, context) {
  const state = context?.getEditorState?.()
  if (!task || context !== task.context || task.editorId !== state?.editorId ||
    task.documentId !== state?.documentId || (context && !state?.active)) return '文章已切换，请重新发起生成'
  if (task.snapshot) {
    const error = snapshotError(task.snapshot, state)
    if (error) return error
  }
  if ((context?.getFullText?.() || '') !== task.payload.content ||
    (context?.getTitle?.() || '') !== task.payload.title) return '文章内容或标题已变化，请使用快捷操作重新生成'
  return ''
}

export function titleFillError(title, expected, current) {
  if (typeof title !== 'string' || !title.trim() || [...title.trim()].length > 120 || /[\r\n<>]/.test(title)) return '候选标题无效'
  if (current !== expected) return '标题已被编辑，请保留当前标题并重新生成候选'
  return ''
}
