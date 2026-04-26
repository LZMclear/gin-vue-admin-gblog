export const normalizePageQuery = (query = {}) => {
  const page = query.page ?? query.pageNum ?? 1
  const normalized = {
    ...query,
    page,
    pageNum: query.pageNum ?? page,
    pageSize: query.pageSize ?? 10
  }

  const dateRange = Array.isArray(query.date) ? query.date : String(query.date || '').split(',')
  if (dateRange.length === 2 && dateRange[0] && dateRange[1]) {
    normalized.startDate = dateRange[0]
    normalized.endDate = dateRange[1]
    delete normalized.date
  }

  return normalized
}

export const asPageResult = (res, list = []) => ({
  ...res,
  data: {
    list,
    total: Array.isArray(list) ? list.length : 0,
    page: 1,
    pageSize: Array.isArray(list) ? list.length : 0
  }
})

export const mapCategoryIn = (item = {}) => ({
  ...item,
  name: item.name ?? item.categoryName
})

export const mapCategoryOut = (form = {}) => ({
  ...form,
  categoryName: form.categoryName ?? form.name
})

export const mapTagIn = (item = {}) => ({
  ...item,
  name: item.name ?? item.tagName
})

export const mapTagOut = (form = {}) => ({
  ...form,
  tagName: form.tagName ?? form.name
})

export const mapArticleOut = (form = {}) => ({
  ...form,
  isPublished: form.isPublished ?? form.published ?? false,
  isRecommend: form.isRecommend ?? form.recommend ?? false,
  isAppreciation: form.isAppreciation ?? form.appreciation ?? false,
  isCommentEnabled: form.isCommentEnabled ?? form.commentEnabled ?? false,
  isTop: form.isTop ?? form.top ?? false,
  categoryId: form.categoryId ?? form.cate,
  tagList: form.tagList ?? []
})

export const mapToggleValue = (id, value) => ({
  id,
  value
})
