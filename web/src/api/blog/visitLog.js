import service from '@/utils/request'
import { normalizePageQuery } from './_helpers'

export function getVisitLogList(queryInfo) {
  return service({
    url: '/admin/visitLogs',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  })
}

export function deleteVisitLogById(id) {
  return service({
    url: '/admin/visitLog',
    method: 'DELETE',
    data: { id }
  })
}

export function deleteVisitLogsByIds(ids) {
  return service({
    url: '/admin/visitLogs',
    method: 'DELETE',
    data: { ids }
  })
}
