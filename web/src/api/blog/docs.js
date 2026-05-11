import service from '@/utils/request'

export function syncDocs() {
	return service({
		url: '/admin/docs/sync',
		method: 'POST'
	})
}
