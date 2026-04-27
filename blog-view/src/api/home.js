import axios from '@/plugins/axios'

export function getBlogList(page) {
	return axios({
		url: 'blogs',
		method: 'GET',
		params: {
			page
		}
	})
}
