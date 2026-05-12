<template>
	<section class="latest-card">
		<header class="latest-header">
			<h2>最新文章</h2>
		</header>
		<BlogItem :blogList="blogList" v-viewer/>
		<Pagination :getBlogList="getBlogList" :totalPage="totalPage"/>
	</section>
</template>

<script>
	import BlogItem from "./BlogItem";
	import Pagination from "./Pagination";

	export default {
		name: "BlogList",
		components: {Pagination, BlogItem},
		props: {
			getBlogList: {
				type: Function,
				required: true
			},
			blogList: {
				type: Array,
				required: true
			},
			totalPage: {
				type: Number,
				required: true
			}
		}
	}
</script>

<style scoped>
	.latest-card {
		padding: 24px;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		background: rgba(255, 255, 255, .9);
		box-shadow: 0 18px 45px rgba(15, 23, 42, .08);
	}

	.latest-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 18px;
	}

	.latest-header h2 {
		margin: 0;
		color: #111827;
		font-size: 22px;
		line-height: 1.2;
		letter-spacing: 0;
	}

	@media only screen and (max-width: 640px) {
		.latest-card {
			padding: 16px;
		}
	}
</style>
