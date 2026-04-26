<template>
	<div>
		<div class="ui top segment" style="text-align: center">
			<h2 class="m-text-500">分类 {{ categoryName }} 下的文章</h2>
		</div>
		<BlogList :getBlogList="getBlogList" :blogList="blogList" :totalPage="totalPage"/>
	</div>
</template>

<script>
	import BlogList from "@/components/blog/BlogList";
	import {getBlogListByCategoryName} from "@/api/category";
	import {getTotalPage, isSuccess, normalizeBlogs} from "@/util/gvaResponse";

	export default {
		name: "Category",
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
				if (this.$route.name === 'category') {
					this.getBlogList()
				}
			}
		},
		created() {
			this.getBlogList()
		},
		computed: {
			categoryName() {
				return this.$route.params.name
			}
		},
		methods: {
			getBlogList(pageNum) {
				getBlogListByCategoryName(this.categoryName, pageNum).then(res => {
					if (isSuccess(res)) {
						if (!res.data || !Array.isArray(res.data.list)) {
							this.blogList = []
							this.totalPage = 0
							this.msgInfo('分类文章列表接口暂未对接')
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
