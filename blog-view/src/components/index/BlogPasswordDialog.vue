<template>
	<!--私密文章密码对话框-->
	<el-dialog title="请输入受保护文章密码" width="30%" :visible.sync="blogPasswordDialogVisible"
	           :lock-scroll="false" :before-close="blogPasswordDialogClosed" @opened="focusPasswordInput">
		<!--内容主体-->
		<el-form :model="blogPasswordForm" :rules="formRules" ref="formRef" label-width="80px">
			<el-form-item label="密码" prop="password">
				<el-input
					ref="passwordInput"
					v-model="blogPasswordForm.password"
					show-password
					@keyup.native.enter="submitBlogPassword"
				></el-input>
			</el-form-item>
		</el-form>
		<!--底部-->
		<span slot="footer">
			<el-button @click="blogPasswordDialogClosed">取 消</el-button>
			<el-button type="primary" @click="submitBlogPassword">确 定</el-button>
		</span>
	</el-dialog>
</template>

<script>
	import {mapState} from "vuex";
	import {SET_BLOG_PASSWORD_DIALOG_VISIBLE} from "../../store/mutations-types";
	import {checkBlogPassword} from "@/api/blog";
	import {isSuccess} from "@/util/gvaResponse";

	export default {
		name: "BlogPasswordDialog",
		computed: {
			...mapState(['blogPasswordDialogVisible', 'blogPasswordForm'])
		},
		data() {
			return {
				formRules: {
					password: [{required: true, message: '请输入密码', trigger: 'change'}]
				}
			}
		},
		methods: {
			focusPasswordInput() {
				this.$nextTick(() => {
					if (this.$refs.passwordInput) {
						this.$refs.passwordInput.focus()
					}
				})
			},
			blogPasswordDialogClosed() {
				this.$refs.formRef.resetFields()
				this.$store.commit(SET_BLOG_PASSWORD_DIALOG_VISIBLE, false)
			},
			submitBlogPassword() {
				this.$refs.formRef.validate(valid => {
					if (valid) {
						checkBlogPassword(this.blogPasswordForm).then(res => {
							if (isSuccess(res)) {
								this.msgSuccess(res.msg)
								window.localStorage.setItem(`blog${this.blogPasswordForm.blogId}`, res.data)
								this.$router.push({name: 'blog', params: {id: this.blogPasswordForm.blogId}}).catch(err => {
									if (err.name !== 'NavigationDuplicated') {
										throw err
									}
								})
								this.blogPasswordDialogClosed()
							} else {
								this.msgError(res.msg)
							}
						}).catch(() => {
							this.msgError("请求失败")
						})
					}
				})
			}
		}
	}
</script>

<style scoped>

</style>
