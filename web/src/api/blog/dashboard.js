import service from '@/utils/request'

export function getDashboard() {
	return service({
		url: '/admin/dashboard',
		method: 'GET'
	})
}