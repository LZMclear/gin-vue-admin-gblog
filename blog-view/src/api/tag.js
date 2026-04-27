import axios from '@/plugins/axios'

export function getBlogListByTagName(tagName, pageNum) {
	return axios({
		url: 'tag/blogs',
		method: 'GET',
		params: {
			tagName,
			page: pageNum || 1
		}
	})
}
