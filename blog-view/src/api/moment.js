import axios from '@/plugins/axios'

export function getMomentListByPageNum(token, page) {
	return axios({
		url: 'moments',
		method: 'GET',
		headers: {
			Authorization: token,
		},
		params: {
			page
		}
	})
}

export function likeMoment(id) {
	return axios({
		url: `moment/like/${id}`,
		method: 'POST',
	})
}
