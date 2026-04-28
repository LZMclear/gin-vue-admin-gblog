import axios from '@/plugins/axios'

export function getBlogListByCategoryName(categoryName, pageNum) {
	return axios({
		url: 'category/blogs',
		method: 'GET',
		params: {
			categoryName,
			page: pageNum || 1
		}
	})
}
