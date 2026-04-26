import service from '@/utils/request'
import { asPageResult, normalizePageQuery } from './_helpers'

export function getData(queryInfo) {
  return service({
    url: '/admin/categories',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  }).then(res => asPageResult(res, res.data || []))
}

export function addCategory(form) {
  return service({
    url: '/admin/category',
    method: 'POST',
    data: form
  })
}

export function editCategory(form) {
  return service({
    url: '/admin/category',
    method: 'PUT',
    data: form
  })
}

export function deleteCategoryById(id) {
  return service({
    url: '/admin/category',
    method: 'DELETE',
    data: { id }
  })
}
