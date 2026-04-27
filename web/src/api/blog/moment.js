import service from '@/utils/request'
import { mapToggleValue, normalizePageQuery } from './_helpers'

export function getMomentListByQuery(queryInfo) {
  return service({
    url: '/admin/moments',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  }).then(res => ({
    ...res,
    data: {
      ...(res.data || {}),
      list: res.data?.list || []
    }
  }))
}

export function updatePublished(id, published) {
  return service({
    url: '/admin/moment/published',
    method: 'PUT',
    params: mapToggleValue(id, published)
  })
}

export function getMomentById(id) {
  return service({
    url: '/admin/moment',
    method: 'GET',
    params: { id }
  })
}

export function deleteMomentById(id) {
  return service({
    url: '/admin/moment',
    method: 'DELETE',
    data: { id }
  })
}

export function saveMoment(moment) {
  return service({
    url: '/admin/moment',
    method: 'POST',
    data: moment
  })
}

export function updateMoment(moment) {
  return service({
    url: '/admin/moment',
    method: 'PUT',
    data: moment
  })
}
