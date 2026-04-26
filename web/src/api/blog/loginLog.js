import service from '@/utils/request'
import { normalizePageQuery } from './_helpers'

export function getLoginLogList(queryInfo) {
  return service({
    url: '/admin/loginLogs',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  })
}

export function deleteLoginLogById(id) {
  return service({
    url: '/admin/loginLog',
    method: 'DELETE',
    data: { id }
  })
}
