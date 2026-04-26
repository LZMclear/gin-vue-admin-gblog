import service from '@/utils/request'
import {
  mapArticleOut,
  mapCategoryIn,
  mapTagIn,
  mapToggleValue,
  normalizePageQuery
} from './_helpers'

const mapArticleIn = (item = {}) => ({
  ...item,
  category: item.category ? mapCategoryIn(item.category) : item.category,
  tags: Array.isArray(item.tags) ? item.tags.map(mapTagIn) : item.tags,
  published: item.published ?? item.isPublished,
  recommend: item.recommend ?? item.isRecommend,
  appreciation: item.appreciation ?? item.isAppreciation,
  commentEnabled: item.commentEnabled ?? item.isCommentEnabled,
  top: item.top ?? item.isTop
})

export function getDataByQuery(queryInfo) {
  return Promise.all([
    service({
      url: '/admin/blogs',
      method: 'GET',
      params: normalizePageQuery(queryInfo)
    }),
    getCategoryAndTag()
  ]).then(([res, metaRes]) => ({
    ...res,
    data: {
      blogs: {
        ...(res.data || {}),
        list: (res.data?.list || []).map(mapArticleIn)
      },
      categories: metaRes.data.categories
    }
  }))
}

export function deleteBlogById(id) {
  return service({
    url: '/admin/blog',
    method: 'DELETE',
    data: { id }
  })
}

export function getCategoryAndTag() {
  return service({
    url: '/admin/categoryAndTag',
    method: 'GET'
  }).then(res => ({
    ...res,
    data: {
      categories: (res.data?.categories || []).map(mapCategoryIn),
      tags: (res.data?.tags || []).map(mapTagIn)
    }
  }))
}

export function saveBlog(blog) {
  return service({
    url: '/admin/blog',
    method: 'POST',
    data: mapArticleOut(blog)
  })
}

export function updateTop(id, top) {
  return service({
    url: '/admin/blog/top',
    method: 'PUT',
    params: mapToggleValue(id, top)
  })
}

export function updateRecommend(id, recommend) {
  return service({
    url: '/admin/blog/recommend',
    method: 'PUT',
    params: mapToggleValue(id, recommend)
  })
}

export function updateVisibility(id, form) {
  return service({
    url: `/admin/blog/${id}/visibility`,
    method: 'PUT',
    data: {
      appreciation: form.appreciation,
      recommend: form.recommend,
      commentEnabled: form.commentEnabled,
      top: form.top,
      published: form.published,
      password: form.password
    }
  })
}

export function getBlogById(id) {
  return service({
    url: '/admin/blog',
    method: 'GET',
    params: { id }
  }).then(res => ({
    ...res,
    data: mapArticleIn(res.data)
  }))
}

export function updateBlog(blog) {
  return service({
    url: '/admin/blog',
    method: 'PUT',
    data: mapArticleOut(blog)
  })
}
