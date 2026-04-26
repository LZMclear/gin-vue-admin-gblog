import service from '@/utils/request'
import { mapToggleValue, normalizePageQuery } from './_helpers'

const mapCommentIn = (item = {}) => ({
  ...item,
  published: item.published ?? item.isPublished,
  adminComment: item.adminComment ?? item.isAdminComment,
  notice: item.notice ?? item.isNotice
})

const mapCommentOut = (item = {}) => ({
  ...item,
  isPublished: item.isPublished ?? item.published ?? false,
  isAdminComment: item.isAdminComment ?? item.adminComment ?? false,
  isNotice: item.isNotice ?? item.notice ?? false
})

export function getCommentListByQuery(queryInfo) {
  return service({
    url: '/admin/comments',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
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
    data: mapCommentOut(form)
  })
}
