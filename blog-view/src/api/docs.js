import axios from '@/plugins/axios'

export function getDocsTree() {
	return axios({
		url: 'docs/tree',
		method: 'GET'
	})
}

export function getDocContent(path) {
	return axios({
		url: 'docs/content',
		method: 'GET',
		params: {
			path
		}
	})
}
