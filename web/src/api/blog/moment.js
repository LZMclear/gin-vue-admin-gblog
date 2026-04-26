import service from '@/utils/request'
import { mapToggleValue, normalizePageQuery } from './_helpers'

const mapMomentIn = (item = {}) => ({
  ...item,
  published: item.published ?? item.isPublished
})

const mapMomentOut = (item = {}) => ({
  ...item,
  isPublished: item.isPublished ?? item.published ?? false
})

export function getMomentListByQuery(queryInfo) {
  return service({
    url: '/admin/moments',
    method: 'GET',
    params: normalizePageQuery(queryInfo)
  }).then(res => ({
    ...res,
    data: {
      ...(res.data || {}),
      list: (res.data?.list || []).map(mapMomentIn)
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
  }).then(res => ({
    ...res,
    data: mapMomentIn(res.data)
  }))
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
    data: mapMomentOut(moment)
  })
}

export function updateMoment(moment) {
  return service({
    url: '/admin/moment',
    method: 'PUT',
    data: mapMomentOut(moment)
  })
}
