import service from '@/utils/request'
import { normalizePageQuery } from './_helpers'

export function getOperationLogList(queryInfo) {
  return service({
    url: '/admin/operationLogs',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  })
}

export function deleteOperationLogById(id) {
  return service({
    url: '/admin/operationLog',
    method: 'DELETE',
    data: { id }
  })
}
