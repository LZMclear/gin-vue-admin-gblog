import service from '@/utils/request'
import { asPageResult, mapCategoryIn, mapCategoryOut, normalizePageQuery } from './_helpers'

export function getData(queryInfo) {
  return service({
    url: '/admin/categories',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  }).then(res => asPageResult(res, (res.data || []).map(mapCategoryIn)))
}

export function addCategory(form) {
  return service({
    url: '/admin/category',
    method: 'POST',
    data: mapCategoryOut(form)
  })
}

export function editCategory(form) {
  return service({
    url: '/admin/category',
    method: 'PUT',
    data: mapCategoryOut(form)
  })
}

export function deleteCategoryById(id) {
  return service({
    url: '/admin/category',
    method: 'DELETE',
    data: { id }
  })
}
