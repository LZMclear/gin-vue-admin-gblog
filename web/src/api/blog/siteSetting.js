import service from '@/utils/request'

export function getSiteSettingData() {
	return service({
		url: '/admin/siteSettings',
		method: 'GET'
	})
}

export function update(settings, deleteIds) {
	return service({
		url: '/admin/siteSettings',
		method: 'POST',
		data: {
			settings,
			deleteIds
		}
	})
}

export function getWebTitleSuffix() {
	return service({
		url: '/admin/webTitleSuffix',
		method: 'GET'
	})
}