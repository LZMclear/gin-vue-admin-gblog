const key = value => String(value ?? '').trim().toLowerCase()

// 仅生成表单变更，全部校验成功后再由页面写入；已有项绝不降级为新建字符串。
export function buildSuggestionPatch(form, suggestion, categories, tags) {
  const patch = {}
  const names = new Map(tags.map(tag => [key(tag.tagName), tag]))
  if (suggestion.category) {
    const category = categories.find(item => key(item.categoryName) === key(suggestion.category))
    if (!category || (suggestion.categoryId && category.id !== suggestion.categoryId)) {
      return { ok: false, message: '推荐分类已不在当前目录中，请刷新后重新生成' }
    }
    patch.cate = category.id
  }
  const additions = []
  for (const [index, name] of (suggestion.tags || []).entries()) {
    const tag = names.get(key(name))
    if (!tag || (suggestion.tagIds?.[index] && tag.id !== suggestion.tagIds[index])) {
      return { ok: false, message: '推荐标签已不在当前目录中，请刷新后重新生成' }
    }
    additions.push(tag.id)
  }
  for (const name of suggestion.newTags || []) {
    if (typeof name !== 'string' || !name.trim() || [...name.trim()].length > 32 || /[<>\p{Cc}]/u.test(name)) {
      return { ok: false, message: '新标签名称无效' }
    }
    additions.push(names.get(key(name))?.id ?? name.trim())
  }
  if (additions.length) {
    const merged = [...(form.tagList || []), ...additions].map(value =>
      typeof value === 'string' ? names.get(key(value))?.id ?? value : value)
    const seen = new Set()
    patch.tagList = merged.filter(value => {
      const identity = typeof value === 'number' ? `id:${value}` : `name:${key(value)}`
      if (seen.has(identity)) return false
      seen.add(identity)
      return true
    })
  }
  if (!Object.keys(patch).length) return { ok: false, message: '请先选择要回填的分类或标签' }
  return { ok: true, patch }
}

export function cursorContext(content, offset) {
  const position = Number.isInteger(offset) ? Math.max(0, Math.min(offset, content.length)) : content.length
  return content.slice(0, position)
}
