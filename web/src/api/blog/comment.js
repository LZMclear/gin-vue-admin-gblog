import service from '@/utils/request'
import { mapToggleValue } from './_helpers'

const mapCommentIn = (item = {}) => ({
  ...item,
  replyComments: item.replyComments || []
})

const normalizeCommentQuery = (query = {}) => {
  const normalized = {
    pageNum: query.pageNum ?? 1,
    pageSize: query.pageSize ?? 10
  }
  if (query.page !== null && query.page !== undefined && query.page !== '') {
    normalized.page = query.page
  }
  if (query.blogId !== null && query.blogId !== undefined && query.blogId !== '') {
    normalized.blogId = query.blogId
  }
  return normalized
}

export function getCommentListByQuery(queryInfo) {
  return service({
    url: '/admin/comments',
    method: 'GET',
    params: normalizeCommentQuery(queryInfo)
  }).then(res => ({
    ...res,
    data: {
      ...(res.data || {}),
      list: (res.data?.list || []).map(mapCommentIn)
    }
  }))
}

export function getBlogList() {
  return service({
    url: '/admin/blogIdAndTitle',
    method: 'GET'
  })
}

export function updatePublished(id, published) {
  return service({
    url: '/admin/comment/published',
    method: 'PUT',
    params: mapToggleValue(id, published)
  })
}

export function updateNotice(id, notice) {
  return service({
    url: '/admin/comment/notice',
    method: 'PUT',
    params: mapToggleValue(id, notice)
  })
}

export function deleteCommentById(id) {
  return service({
    url: '/admin/comment',
    method: 'DELETE',
    data: { id }
  })
}

export function editComment(form) {
  return service({
    url: '/admin/comment',
    method: 'PUT',
    data: form
  })
}
