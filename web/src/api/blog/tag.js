import service from '@/utils/request'
import { asPageResult, normalizePageQuery } from './_helpers'

export function getData(queryInfo) {
  return service({
    url: '/admin/tags',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  }).then(res => asPageResult(res, res.data || []))
}

export function addTag(form) {
  return service({
    url: '/admin/tag',
    method: 'POST',
    data: form
  })
}

export function editTag(form) {
  return service({
    url: '/admin/tag',
    method: 'PUT',
    data: form
  })
}

export function deleteTagById(id) {
  return service({
    url: '/admin/tag',
    method: 'DELETE',
    data: { id }
  })
}
