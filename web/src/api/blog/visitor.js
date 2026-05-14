import service from '@/utils/request'
import { normalizePageQuery } from './_helpers'

export function getVisitorList(queryInfo) {
  return service({
    url: '/admin/visitors',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  })
}

export function deleteVisitor(id) {
  return service({
    url: '/admin/visitor',
    method: 'DELETE',
    data: { id }
  })
}

export function deleteVisitorsByIds(ids) {
  return service({
    url: '/admin/visitors',
    method: 'DELETE',
    data: { ids }
  })
}
