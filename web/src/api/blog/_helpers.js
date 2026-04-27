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

const toNumber = (value, fallback = 0) => {
  if (value === '' || value === null || value === undefined) {
    return fallback
  }
  const numberValue = Number(value)
  return Number.isFinite(numberValue) ? numberValue : fallback
}

export const mapArticleOut = (form = {}) => ({
  id: toNumber(form.id),
  title: form.title,
  firstPicture: form.firstPicture,
  content: form.content,
  description: form.description,
  words: toNumber(form.words),
  readTime: toNumber(form.readTime),
  views: toNumber(form.views),
  password: form.password,
  cate: form.cate,
  isPublished: form.published,
  isRecommend: form.recommend,
  isAppreciation: form.appreciation,
  isCommentEnabled: form.commentEnabled,
  isTop: form.top,
  categoryId: toNumber(form.categoryId ?? form.cate),
  tagList: form.tagList ?? []
})

export const mapToggleValue = (id, value) => ({
  id,
  value
})
