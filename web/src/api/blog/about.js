import service from '@/utils/request'

export function getAbout() {
  return service({
    url: '/admin/about',
    method: 'GET'
  }).then(res => {
    const values = {}
    ;(res.data || []).forEach(item => {
      values[item.nameEn] = item.value
    })
    return {
      ...res,
      data: {
        title: values.title || '',
        musicId: values.musicId || '',
        content: values.content || '',
        commentEnabled: values.commentEnabled || 'true'
      }
    }
  })
}

export function updateAbout(form) {
  return service({
    url: '/admin/about',
    method: 'PUT',
    data: {
      values: {
        title: String(form.title ?? ''),
        musicId: String(form.musicId ?? ''),
        content: String(form.content ?? ''),
        commentEnabled: String(Boolean(form.commentEnabled))
      }
    }
  })
}
