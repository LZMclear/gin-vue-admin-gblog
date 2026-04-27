import {
	SAVE_COMMENT_RESULT,
	SET_PARENT_COMMENT_ID,
	RESET_COMMENT_FORM,
	SET_BLOG_PASSWORD_DIALOG_VISIBLE,
	SET_BLOG_PASSWORD_FORM
} from "./mutations-types";

import {getCommentListByQuery, submitComment} from "@/api/comment";
import {Message, Notification} from "element-ui";
import router from "../router";
import sanitizeHtml from 'sanitize-html'
import {isSuccess, normalizeCommentPage} from '@/util/gvaResponse'

function pushRoute(location) {
	return router.push(location).catch(err => {
		if (err && err.name !== 'NavigationDuplicated') {
			throw err
		}
	})
}

export default {
	getCommentList({commit, rootState}) {
		//密码保护的文章，需要发送密码验证通过后保存在localStorage的Token
		const blogToken = window.localStorage.getItem(`blog${rootState.commentQuery.blogId}`)
		//如果有则发送博主身份Token
		const adminToken = window.localStorage.getItem('adminToken')
		const token = adminToken ? adminToken : (blogToken ? blogToken : '')

		getCommentListByQuery(token, rootState.commentQuery).then(res => {
			if (isSuccess(res)) {
				let sanitizeHtmlConfig = {
					allowedTags: [],
					allowedAttributes: false,
					disallowedTagsMode: 'recursiveEscape'
				}
				const data = normalizeCommentPage(res.data)
				data.comments.list.forEach(comment => {
					//转义评论中的html
					comment.content = sanitizeHtml(comment.content, sanitizeHtmlConfig)
					//查找评论中是否有表情
					comment.replyComments = comment.replyComments || []
					comment.replyComments.forEach(comment => {
						//转义评论中的html
						comment.content = sanitizeHtml(comment.content, sanitizeHtmlConfig)
						//查找评论中是否有表情
					})
				})
				commit(SAVE_COMMENT_RESULT, data)
			}
		}).catch(() => {
			Message.error("请求失败")
		})
	},
	submitCommentForm({rootState, dispatch, commit}, token) {
		let form = {...rootState.commentForm}
		form.page = rootState.commentQuery.page
		form.blogId = rootState.commentQuery.blogId
		form.parentCommentId = rootState.parentCommentId
		form.isNotice = form.notice
		submitComment(token, form).then(res => {
			if (isSuccess(res)) {
				Notification({
					title: res.msg,
					type: 'success'
				})
				commit(SET_PARENT_COMMENT_ID, -1)
				commit(RESET_COMMENT_FORM)
				dispatch('getCommentList')
			} else {
				Notification({
					title: '评论失败',
					message: res.msg,
					type: 'error'
				})
			}
		}).catch(() => {
			Notification({
				title: '评论失败',
				message: '异常错误',
				type: 'error'
			})
		})
	},
	goBlogPage({commit}, blog) {
		if (blog.privacy) {
			const adminToken = window.localStorage.getItem('adminToken')
			const blogToken = window.localStorage.getItem(`blog${blog.id}`)
			//对于密码保护文章，博主身份Token和经过密码验证后的Token都可以跳转路由，再由后端验证Token有效性
			if (adminToken || blogToken) {
				return pushRoute({name: 'blog', params: {id: blog.id}})
			}
			commit(SET_BLOG_PASSWORD_FORM, {blogId: blog.id, password: ''})
			commit(SET_BLOG_PASSWORD_DIALOG_VISIBLE, true)
		} else {
			pushRoute({name: 'blog', params: {id: blog.id}})
		}
	},
}
