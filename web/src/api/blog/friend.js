import service from '@/utils/request'
import { mapToggleValue, normalizePageQuery } from './_helpers'

export function getFriendsByQuery(queryInfo) {
  return service({
    url: '/admin/friends',
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
    url: '/admin/friend/published',
    method: 'PUT',
    params: mapToggleValue(id, published)
  })
}

export function saveFriend(form) {
  return service({
    url: '/admin/friend',
    method: 'POST',
    data: form
  })
}

export function updateFriend(form) {
  return service({
    url: '/admin/friend',
    method: 'PUT',
    data: form
  })
}

export function deleteFriendById(id) {
  return service({
    url: '/admin/friend',
    method: 'DELETE',
    data: { id }
  })
}

export function getFriendInfo() {
  return service({
    url: '/admin/friendInfo',
    method: 'GET'
  })
}

export function updateCommentEnabled(commentEnabled) {
  return service({
    url: '/admin/friendInfo/commentEnabled',
    method: 'PUT',
    params: { id: 0, value: commentEnabled }
  })
}

export function updateContent(content) {
  return service({
    url: '/admin/friendInfo/content',
    method: 'PUT',
    data: { content }
  })
}
