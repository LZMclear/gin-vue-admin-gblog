<template>
	<!-- 评论输入表单 -->
	<div class="form">
		<h3>
			发表评论
			<el-button class="m-small" size="mini" type="primary" @click="$store.commit(SET_PARENT_COMMENT_ID, -1)" v-show="parentCommentId!==-1">取消回复</el-button>
		</h3>
		<el-form :inline="true" :model="commentForm" :rules="formRules" ref="formRef" size="small">
			<el-input :class="'textarea'" type="textarea" :rows="5" v-model="commentForm.content" placeholder="评论千万条，友善第一条"
			          maxlength="250" show-word-limit :validate-event="false"></el-input>
			<div class="el-form-item el-form-item--small emoji">
				<button type="button" class="emoji-trigger" @click="showEmojiBox">😀</button>
				<div class="mask" v-show="emojiShow" @click="hideEmojiBox"></div>
				<div class="emoji-box" v-show="emojiShow">
					<div class="emoji-title">
						<span>Emoji</span>
					</div>
					<div class="emoji-wrap">
						<button type="button" class="emoji-list" v-for="emoji in emojis" :key="emoji" @click="insertEmoji(emoji)">
							{{ emoji }}
						</button>
					</div>
				</div>
			</div>


			<!-- 评论 -->
			<el-form-item prop="nickname">
				<el-popover ref="nicknamePopover" placement="bottom" trigger="focus" content="输入QQ号将自动拉取昵称和头像"></el-popover>
				<el-input v-model="commentForm.nickname" placeholder="昵称（必填）" :validate-event="false" v-popover:nicknamePopover @blur="fillQQInfo">
					<i slot="prefix" class="el-input__icon el-icon-user"></i>
				</el-input>
			</el-form-item>
			<el-form-item prop="email">
				<el-popover ref="emailPopover" placement="bottom" trigger="focus" content="用于接收回复邮件"></el-popover>
				<el-input v-model="commentForm.email" placeholder="邮箱（必填）" :validate-event="false" v-popover:emailPopover>
					<i slot="prefix" class="el-input__icon el-icon-message"></i>
				</el-input>
			</el-form-item>
			<el-form-item prop="website">
				<el-popover ref="websitePopover" placement="bottom" trigger="focus" content="可以让我参观一下吗😊"></el-popover>
				<el-input v-model="commentForm.website" placeholder="https://（可选）" :validate-event="false" v-popover:websitePopover>
					<i slot="prefix" class="el-input__icon el-icon-map-location"></i>
				</el-input>
			</el-form-item>
			<el-form-item label="订阅回复">
				<el-switch v-model="commentForm.isNotice"></el-switch>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" size="medium" v-throttle="[postForm,`click`,3000]">发表评论</el-button>
			</el-form-item>
		</el-form>
	</div>
</template>

<script>
	import {mapState} from 'vuex'
	import {checkEmail, checkUrl} from "@/common/reg";
	import {SET_PARENT_COMMENT_ID} from "@/store/mutations-types";

	const validateWebsite = (rule, value, callback) => {
		if (value) {
			return checkUrl(rule, value, callback)
		}
		callback()
	}
	export default {
		name: "CommentForm",
		computed: {
			...mapState(['parentCommentId', 'commentForm', 'commentQuery'])
		},
		data() {
			return {
				SET_PARENT_COMMENT_ID,
				formRules: {
					nickname: [
						{required: true, message: '请输入评论昵称'},
						{max: 18, message: '昵称不可多于15个字符'}
					],
					email: [
						{required: true, message: '请输入评论邮箱'},
						{validator: checkEmail}
					],
					website: [
						{required: false},
						{validator: validateWebsite}
					]
				},
				emojiShow: false,
				emojis: [
					'😀', '😄', '😂', '🤣', '😊', '😍', '😘', '😎',
					'🤔', '😮', '😅', '😭', '😡', '😴', '🙄', '😇',
					'👍', '👎', '👏', '🙏', '💪', '🤝', '👌', '✌️',
					'❤️', '💔', '🔥', '🎉', '✨', '🌹', '🍉', '☕'
				],
				textarea: null,
				start: 0,
				end: 0,
			}
		},
		mounted() {
			this.textarea = document.querySelector('.el-form textarea')
		},
		methods: {
			showEmojiBox() {
				this.start = this.textarea.selectionStart
				this.end = this.textarea.selectionEnd
				this.textarea.focus()
				this.textarea.setSelectionRange(this.start, this.end)
				this.emojiShow = !this.emojiShow
			},
			insertEmoji(name) {
				let str = this.commentForm.content
				this.commentForm.content = str.substring(0, this.start) + name + str.substring(this.end)
				this.start += name.length
				this.end = this.start
				this.textarea.focus()
				this.$nextTick(() => {
					this.textarea.setSelectionRange(this.start, this.end)
				})
			},
			hideEmojiBox() {
				this.emojiShow = false
				this.textarea.focus()
				this.textarea.setSelectionRange(this.start, this.end)
			},
			fillQQInfo() {
				const nickname = String(this.commentForm.nickname || '').trim()
				if (!nickname) {
					return
				}
				if (!/^[1-9][0-9]{4,11}$/.test(nickname)) {
					this.commentForm.qq = ''
					if (!this.commentForm.avatar || this.commentForm.avatar.includes('qlogo.cn')) {
						this.commentForm.avatar = this.randomCommentAvatar(nickname)
					}
					return
				}
				const qq = nickname
				this.commentForm.qq = qq
				this.commentForm.avatar = `https://q1.qlogo.cn/g?b=qq&nk=${qq}&s=100`
				if (!this.commentForm.email) {
					this.commentForm.email = `${qq}@qq.com`
				}
				fetch(`https://api.uomg.com/api/qq.info?qq=${qq}`)
					.then(res => res.json())
					.then(data => {
						if (data && data.code === 1 && data.name) {
							this.commentForm.nickname = data.name
						}
						if (data && data.qlogo) {
							this.commentForm.avatar = data.qlogo
						}
					})
					.catch(() => {})
			},
			randomCommentAvatar(seed) {
				const text = String(seed || '')
				let hash = 0
				for (let i = 0; i < text.length; i++) {
					hash = (hash * 31 + text.charCodeAt(i)) >>> 0
				}
				return `/img/comment-avatar/${hash % 6 + 1}.jpg`
			},
			postForm() {
				const adminToken = window.localStorage.getItem('adminToken')
				if (adminToken) {
					//博主登录后，localStorage中会存储token，在后端设置属性，可以不校验昵称、邮箱
					if (this.commentForm.content === '' || this.commentForm.content.length > 250) {
						return this.$notify({
							title: '评论失败',
							message: '评论内容有误',
							type: 'warning'
						})
					} else {
						return this.$store.dispatch('submitCommentForm', adminToken)
					}
				}
				const blogToken = window.localStorage.getItem(`blog${this.commentQuery.blogId}`)
				this.$refs.formRef.validate(valid => {
					if (!valid || this.commentForm.content === '' || this.commentForm.content.length > 250) {
						this.$notify({
							title: '评论失败',
							message: '请正确填写评论',
							type: 'warning'
						})
					} else {
						this.$store.dispatch('submitCommentForm', blogToken ? blogToken : '')
					}
				})
			}
		}
	}
</script>

<style>
	.form {
		background: #fff;
		position: relative;
	}

	.form h3 {
		margin: 5px;
		font-weight: 500 !important;
	}

	.form .m-small {
		margin-left: 5px;
		padding: 4px 5px;
	}

	.el-form .textarea {
		margin-top: 5px;
		margin-bottom: 15px;
	}

	.el-form textarea {
		padding: 6px 8px;
	}

	.el-form textarea, .el-form input {
		color: black;
	}

	.el-form .el-form-item__label {
		padding-right: 3px;
	}

	.emoji {
		margin-right: 5px;
		position: relative;
		user-select: none;
	}

	.emoji-trigger {
		width: 32px;
		height: 32px;
		padding: 0;
		border: 0;
		background: transparent;
		cursor: pointer;
		font-size: 24px;
		line-height: 32px;
		transition: all 0.3s ease-in-out;
		-webkit-transition: all 0.3s ease-in-out;
		-moz-transition: all 0.3s ease-in-out;
		-o-transition: all 0.3s ease-in-out;
	}

	.emoji-trigger:hover {
		transform: rotate(360deg);
		-webkit-transform: rotate(360deg);
		-moz-transform: rotate(360deg);
		-o-transform: rotate(360deg);
	}

	.emoji-box {
		color: #222;
		overflow: visible;
		background: #fff;
		border: 1px solid #E5E9EF;
		box-shadow: 0 11px 12px 0 rgba(106, 115, 133, 0.3);
		border-radius: 8px;
		width: 340px;
		position: absolute;
		top: 40px;
		z-index: 100;
	}

	.emoji-box * {
		box-sizing: content-box;
	}

	.emoji-box .emoji-title {
		font-size: 12px;
		line-height: 16px;
		margin: 13px 15px 0;
		color: #757575;
	}

	.emoji-box .emoji-wrap {
		margin: 6px 11px 0 11px;
		max-height: 185px;
		overflow: auto;
		word-break: break-word;
	}

	.emoji-box .emoji-wrap .emoji-list {
		width: 36px;
		height: 36px;
		padding: 0;
		border: 0;
		background: transparent;
		color: #111;
		border-radius: 4px;
		transition: background 0.2s;
		display: inline-block;
		outline: 0;
		cursor: pointer;
		font-size: 22px;
		line-height: 36px;
		text-align: center;
	}

	.emoji-box .emoji-wrap .emoji-list:hover {
		background-color: #ddd;
	}

	.mask {
		pointer-events: auto;
		position: fixed;
		z-index: 99;
		top: 0;
		bottom: 0;
		left: 0;
		right: 0;
	}
</style>
