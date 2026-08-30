<template>
	<div class="ui segments m-box site-about">
		<div class="ui secondary segment site-about-title">
			<i class="address book outline icon"></i>
			<span>关于本站</span>
		</div>
		<div class="ui yellow segment site-about-body">

			<div class="site-about-slogan">
				<p>记录生活，记录人生</p>
				<p>热爱技术、写作与思考</p>
				<p>用文字和照片记录平凡点滴美好</p>
			</div>
			<div class="site-about-stats">
				<div class="stat-item">
					<div class="stat-label">文章</div>
					<div class="stat-value">{{ stats.articleCount || 0 }}</div>
				</div>
				<div class="stat-item">
					<div class="stat-label">分类</div>
					<div class="stat-value">{{ stats.categoryCount || 0 }}</div>
				</div>
				<div class="stat-item">
					<div class="stat-label">标签</div>
					<div class="stat-value">{{ stats.tagCount || 0 }}</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
	export default {
		name: "SiteAbout",
		props: {
			stats: {
				type: Object,
				default: () => ({
					articleCount: 0,
					categoryCount: 0,
					tagCount: 0
				})
			}
		}
	}
</script>

<style scoped>
	.site-about {
		border-radius: 8px !important;
		overflow: hidden;
		box-shadow: 0 8px 24px rgba(15, 23, 42, .08) !important;
	}

	.site-about-title {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 12px 14px !important;
		color: #1f2937;
		font-size: 16px;
		font-weight: 700;
	}

	.site-about-title i {
		display: inline-flex !important;
		align-items: center;
		justify-content: center;
		width: 18px;
		height: 18px;
		margin: 0 !important;
		color: #2563eb;
		font-size: 16px;
		line-height: 1 !important;
	}

	.site-about-title span {
		display: inline-flex;
		align-items: center;
		line-height: 1;
	}

	.site-about-body {
		padding: 18px 16px 20px !important;
		background: #fff !important;
	}

	.site-about-copy,
	.site-about-slogan {
		color: #64748b;
		font-size: 14px;
		font-weight: 600;
		line-height: 1.85;
		text-align: left;
	}

	.site-about-copy p,
	.site-about-slogan p {
		margin: 0;
	}

	.site-about-slogan {
		margin-top: 16px;
		text-align: center;
	}

	.site-about-stats {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		margin-top: 18px;
		padding-top: 16px;
		border-top: 1px solid #eef2f7;
	}

	.stat-item {
		text-align: center;
	}

	.stat-item + .stat-item {
		border-left: 1px solid #e5e7eb;
	}

	.stat-label {
		color: #64748b;
		font-size: 13px;
		font-weight: 600;
		line-height: 1.2;
	}

	.stat-value {
		margin-top: 8px;
		color: #1f2937;
		font-size: 18px;
		font-weight: 700;
		line-height: 1.2;
	}
</style>
