import service from '@/utils/request'
import { normalizePageQuery } from './_helpers'

export function getExceptionLogList(queryInfo) {
  return service({
    url: '/admin/exceptionLogs',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  })
}

export function deleteExceptionLogById(id) {
  return service({
    url: '/admin/exceptionLog',
    method: 'DELETE',
    data: { id }
  })
}
