export function isSuccess(res) {
	return res && res.code === 0
}

export function getTotalPage(pageData) {
	if (!pageData) {
		return 0
	}
	if (pageData.totalPage !== undefined) {
		return pageData.totalPage
	}
	const pageSize = pageData.pageSize || 10
	return Math.ceil((pageData.total || 0) / pageSize)
}

export function normalizeCategory(category) {
	return category
}

export function normalizeTag(tag) {
	return tag
}

export function normalizeBlog(blog) {
	if (!blog) {
		return blog
	}
	return {
		...blog,
		day: blog.day || getDay(blog.createTime),
		category: normalizeCategory(blog.category),
		tags: Array.isArray(blog.tags) ? blog.tags.map(normalizeTag) : []
	}
}

export function normalizeBlogs(blogs) {
	return Array.isArray(blogs) ? blogs.map(normalizeBlog) : []
}

export function normalizeSite(data = {}) {
	const siteInfo = {}
	const introduction = {
		avatar: '',
		name: '',
		rollText: [],
		favorites: []
	}
	const badges = []

	;(data.siteSettings || []).forEach(item => {
		const key = item.nameEn
		const value = item.value
		if (!key) {
			return
		}

		if (item.type === 1) {
			siteInfo[key] = parseSettingValue(value)
		} else if (item.type === 2) {
			if (key === 'favorite') {
				introduction.favorites.push(parseSettingValue(value))
			} else if (key === 'rollText') {
				introduction.rollText = parseRollText(value)
			} else {
				introduction[key] = value
			}
		} else if (item.type === 3) {
			badges.push(parseSettingValue(value))
		}
	})

	return {
		siteInfo: data.siteInfo || siteInfo,
		introduction: data.introduction || introduction,
		badges: data.badges || badges,
		categoryList: (data.categoryList || []).map(normalizeCategory),
		tagList: (data.tagList || []).map(normalizeTag),
		newBlogList: normalizeBlogs(data.newBlogList || []),
		randomBlogList: normalizeBlogs(data.randomBlogList || [])
	}
}

export function normalizeAbout(data) {
	if (!Array.isArray(data)) {
		return data || {}
	}
	return data.reduce((result, item) => {
		result[item.nameEn] = item.value
		return result
	}, {})
}

export function normalizeFriendPage(data = {}) {
	return data
}

export function normalizeMoments(moments) {
	return Array.isArray(moments) ? moments : []
}

export function normalizeCommentPage(data = {}) {
	const comments = data.comments || {}
	const list = Array.isArray(comments.list) ? comments.list : []
	const normalizedList = normalizeComments(list)
	return {
		...data,
		comments: {
			...comments,
			list: normalizedList
		}
	}
}

function normalizeComments(list) {
	const commentMap = {}
	const roots = []

	list.forEach(item => {
		const comment = {
			...item,
			replyComments: []
		}
		commentMap[comment.id] = comment
	})

	Object.keys(commentMap).forEach(key => {
		const comment = commentMap[key]
		if (comment.parentCommentId && comment.parentCommentId !== -1 && commentMap[comment.parentCommentId]) {
			commentMap[comment.parentCommentId].replyComments.push(comment)
		} else {
			roots.push(comment)
		}
	})

	return roots
}

function parseSettingValue(value) {
	if (typeof value !== 'string') {
		return value
	}
	try {
		return JSON.parse(value)
	} catch (e) {
		return value
	}
}

function getDay(value) {
	if (!value) {
		return ''
	}
	const date = new Date(value)
	if (Number.isNaN(date.getTime())) {
		return ''
	}
	return String(date.getDate()).padStart(2, '0')
}

function parseRollText(value) {
	if (!value) {
		return []
	}
	try {
		const parsed = JSON.parse(`[${value}]`)
		return Array.isArray(parsed) ? parsed : []
	} catch (e) {
		return value.match(/"([^"]*)"/g)?.map(item => item.slice(1, -1)) || []
	}
}
