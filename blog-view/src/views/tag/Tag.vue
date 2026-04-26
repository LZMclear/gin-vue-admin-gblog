<template>
	<div>
		<div class="ui top segment" style="text-align: center">
			<h2 class="m-text-500">标签 {{ tagName }} 下的文章</h2>
		</div>
		<BlogList :getBlogList="getBlogList" :blogList="blogList" :totalPage="totalPage"/>
	</div>
</template>

<script>
	import BlogList from "@/components/blog/BlogList";
	import {getBlogListByTagName} from "@/api/tag";
	import {getTotalPage, isSuccess, normalizeBlogs} from "@/util/gvaResponse";

	export default {
		name: "Tag",
		components: {BlogList},
		data() {
			return {
				blogList: [],
				totalPage: 0
			}
		},
		watch: {
			//在当前组件被重用时，要重新获取博客列表
			'$route.fullPath'() {
				if (this.$route.name === 'tag') {
					this.getBlogList()
				}
			}
		},
		created() {
			this.getBlogList()
		},
		computed: {
			tagName() {
				return this.$route.params.name
			}
		},
		methods: {
			getBlogList(pageNum) {
				getBlogListByTagName(this.tagName, pageNum).then(res => {
					if (isSuccess(res)) {
						if (!res.data || !Array.isArray(res.data.list)) {
							this.blogList = []
							this.totalPage = 0
							this.msgInfo('标签文章列表接口暂未对接')
							return
						}
						this.blogList = normalizeBlogs(res.data.list)
						this.totalPage = getTotalPage(res.data)
						this.$nextTick(() => {
							Prism.highlightAll()
						})
					} else {
						this.msgError(res.msg)
					}
				}).catch(() => {
					this.msgError("请求失败")
				})
			}
		}
	}
</script>

<style scoped>

</style>
