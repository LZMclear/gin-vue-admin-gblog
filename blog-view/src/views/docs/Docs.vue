<template>
	<div class="docs-page">
		<div class="docs-mobile-picker">
			<el-select :value="activePath" placeholder="选择文档" filterable size="small" @change="selectDocPath">
				<el-option
					v-for="item in flatDocs"
					:key="item.path"
					:label="item.label"
					:value="item.path"
				/>
			</el-select>
		</div>
		<aside class="docs-sidebar m-mobile-hide">
			<div class="docs-panel">
				<div class="docs-panel-title">
					<i class="book icon"></i>
					<span>文档目录</span>
				</div>
				<div v-if="loadingTree" class="docs-empty">目录加载中...</div>
				<div v-else-if="tree.length === 0" class="docs-empty">暂无文档</div>
				<DocTree v-else :nodes="tree" :active-path="activePath" @select="selectDoc"/>
			</div>
		</aside>

		<main class="docs-content">
			<div class="ui padded attached segment docs-article">
				<div v-if="loadingContent" class="docs-empty">文档加载中...</div>
				<div v-else-if="error" class="docs-empty error">{{ error }}</div>
				<template v-else>
					<header v-if="doc.path" class="blog-header">
						<div class="blog-header-top">
							<div class="blog-header-left">
								<span class="docs-path">{{ doc.path }}</span>
							</div>
							<a v-if="doc.editUrl" :href="doc.editUrl" target="_blank" rel="noopener noreferrer" class="header-category">
								<i class="github icon"></i>在 GitHub 上编辑
							</a>
						</div>
						<h1 class="blog-title">{{ docFileTitle }}</h1>
						<div class="blog-meta">
							<span><i class="user outline icon"></i>作者：Gvto</span>
							<span><i class="book icon"></i>字数 {{ docWordCount }}</span>
							<span><i class="clock outline icon"></i>阅读时长 {{ docReadMinutes }} 分钟</span>
							<button class="meta-action" type="button" @click.prevent="bigFontSize=!bigFontSize" title="切换字体大小">
								<i class="font icon"></i>
							</button>
						</div>
					</header>
					<div class="ui middle aligned mobile reversed stackable">
						<div class="ui grid m-margin-lr">
							<div
								class="typo js-toc-content m-padded-tb-small match-braces rainbow-braces"
								v-viewer
								:class="{'m-big-fontsize':bigFontSize}"
								v-html="doc.content"
							></div>
						</div>
					</div>
				</template>
			</div>
		</main>

		<aside class="docs-toc m-mobile-hide">
			<Tocbot
				title="文档目录"
				empty-text="暂无目录"
				content-selector=".docs-page .js-toc-content"
				:refresh-key="doc.path"
				:number-content-headings="true"
			/>
		</aside>
	</div>
</template>

<script>
	import DocTree from '@/components/docs/DocTree'
	import Tocbot from '@/components/sidebar/Tocbot'
	import {getDocContent, getDocsTree} from '@/api/docs'
	import {isSuccess} from '@/util/gvaResponse'

	export default {
		name: 'Docs',
		components: {DocTree, Tocbot},
		data() {
			return {
				tree: [],
				doc: {
					title: '',
					path: '',
					content: '',
					editUrl: ''
				},
				loadingTree: false,
				loadingContent: false,
				error: '',
				bigFontSize: false,
				contentRequestSeq: 0
			}
		},
		computed: {
			activePath() {
				return this.$route.query.path || ''
			},
			flatDocs() {
				const result = []
				const walk = (nodes, parents = []) => {
					nodes.forEach(node => {
						const nextParents = node.type === 'dir' ? [...parents, node.title] : parents
						if (node.type === 'file' && node.path) {
							result.push({
								path: node.path,
								label: [...parents, node.title].join(' / ')
							})
						}
						if (node.children && node.children.length) {
							walk(node.children, nextParents)
						}
					})
				}
				walk(this.tree)
				return result
			},
			docFileTitle() {
				const filename = (this.doc.path || '').split('/').pop() || this.doc.title || 'README'
				return filename.replace(/\.md$/i, '') || 'README'
			},
			docPlainText() {
				const wrapper = document.createElement('div')
				wrapper.innerHTML = this.doc.content || ''
				return (wrapper.textContent || wrapper.innerText || '').trim()
			},
			docWordCount() {
				const text = this.docPlainText
				if (!text) {
					return 0
				}
				const cjkCount = (text.match(/[\u4e00-\u9fa5]/g) || []).length
				const wordCount = (text.replace(/[\u4e00-\u9fa5]/g, ' ').match(/[A-Za-z0-9_]+(?:[-'][A-Za-z0-9_]+)*/g) || []).length
				return cjkCount + wordCount
			},
			docReadMinutes() {
				return Math.max(1, Math.ceil(this.docWordCount / 400))
			}
		},
		watch: {
			'$route.query.path'(path) {
				if (path) {
					this.loadContent(path)
				}
			}
		},
		created() {
			this.loadTree()
		},
		methods: {
			loadTree() {
				this.loadingTree = true
				getDocsTree().then(res => {
					if (!isSuccess(res)) {
						this.error = res.msg || '获取文档目录失败'
						return
					}
					this.tree = Array.isArray(res.data) ? res.data : []
					const currentPath = this.activePath
					const firstPath = currentPath || this.findReadmeDocPath(this.tree) || this.findFirstDocPath(this.tree)
					if (firstPath && firstPath !== currentPath) {
						this.$router.replace({name: 'docs', query: {path: firstPath}})
					} else if (firstPath) {
						this.loadContent(firstPath)
					}
				}).catch(() => {
					this.error = '获取文档目录失败'
				}).finally(() => {
					this.loadingTree = false
				})
			},
			loadContent(path) {
				const requestSeq = ++this.contentRequestSeq
				this.loadingContent = true
				this.error = ''
				getDocContent(path).then(res => {
					if (requestSeq !== this.contentRequestSeq) {
						return
					}
					if (!isSuccess(res)) {
						this.error = res.msg || '获取文档内容失败'
						this.loadingContent = false
						return
					}
					this.doc = res.data || {}
					document.title = `${this.doc.title || '文档'}${this.$store.state.siteInfo.webTitleSuffix || ''}`
					this.loadingContent = false
					this.$nextTick(() => {
						if (window.Prism && typeof window.Prism.highlightAll === 'function') {
							window.Prism.highlightAll()
						}
					})
				}).catch(() => {
					if (requestSeq !== this.contentRequestSeq) {
						return
					}
					this.error = '获取文档内容失败'
					this.loadingContent = false
				})
			},
			selectDoc(node) {
				if (node.path === this.activePath) {
					return
				}
				this.$router.push({name: 'docs', query: {path: node.path}})
			},
			selectDocPath(path) {
				if (!path || path === this.activePath) {
					return
				}
				this.$router.push({name: 'docs', query: {path}})
			},
			findFirstDocPath(nodes) {
				for (const node of nodes) {
					if (node.type === 'file' && node.path) {
						return node.path
					}
					if (node.children && node.children.length) {
						const found = this.findFirstDocPath(node.children)
						if (found) {
							return found
						}
					}
				}
				return ''
			},
			findReadmeDocPath(nodes) {
				for (const node of nodes) {
					if (node.type === 'file' && node.path && /(^|\/)README\.md$/i.test(node.path)) {
						return node.path
					}
					if (node.children && node.children.length) {
						const found = this.findReadmeDocPath(node.children)
						if (found) {
							return found
						}
					}
				}
				return ''
			},
		}
	}
</script>

<style scoped>
	.docs-page {
		display: grid;
		grid-template-columns: 260px minmax(0, 1fr) 240px;
		gap: 20px;
		align-items: start;
	}

	.docs-mobile-picker {
		display: none;
	}

	.docs-sidebar {
		position: sticky;
		top: 62px;
		max-height: calc(100vh - 82px);
		overflow: auto;
		scrollbar-width: thin;
	}

	.docs-toc {
		position: relative;
	}

	.docs-panel {
		background: #fff;
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		padding: 14px;
		box-shadow: 0 8px 24px rgba(15, 23, 42, .04);
	}

	.docs-panel-title {
		display: flex;
		align-items: center;
		gap: 8px;
		padding-bottom: 12px;
		margin-bottom: 10px;
		border-bottom: 1px solid #eef2f7;
		color: #1f2937;
		font-weight: 700;
		font-size: 15px;
	}

	.docs-content {
		min-width: 0;
	}

	.docs-article {
		position: relative;
		min-height: 520px;
		overflow: hidden;
		background: #fff !important;
		border: 1px solid #e5e7eb !important;
		border-radius: 6px !important;
		box-shadow: 0 10px 30px rgba(15, 23, 42, .035);
	}

	.blog-header {
		margin: 0 1rem 24px;
		padding: 0 0 24px;
		border-bottom: 1px solid #e5e7eb;
		background: #fff;
	}

	.blog-header-top {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		margin-bottom: 18px;
	}

	.blog-header-left {
		display: flex;
		min-width: 0;
		align-items: center;
		gap: 8px;
		color: #6b7280;
		font-size: 13px;
		font-weight: 500;
	}

	.docs-path {
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.header-category {
		display: inline-flex;
		flex: 0 0 auto;
		min-height: 30px;
		align-items: center;
		justify-content: center;
		gap: 6px;
		padding: 7px 12px;
		border-radius: 6px;
		background: #eff6ff;
		color: #2563eb;
		font-size: 13px;
		font-weight: 700;
		line-height: 1;
		white-space: nowrap;
	}

	.header-category:hover {
		background: #dbeafe;
		color: #1d4ed8;
	}

	.header-category i,
	.blog-meta i,
	.meta-action i {
		display: inline-flex;
		width: 16px;
		height: 16px;
		align-items: center;
		justify-content: center;
		margin: 0 !important;
		line-height: 1 !important;
	}

	.blog-title {
		margin: 0;
		color: #0f172a;
		font-size: 34px;
		font-weight: 800;
		letter-spacing: 0;
		line-height: 1.25;
	}

	.blog-meta {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 14px;
		margin-top: 20px;
		color: #6b7280;
		font-size: 14px;
		font-weight: 400;
	}

	.blog-meta span {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		line-height: 1;
	}

	.meta-action {
		display: inline-flex;
		width: 28px;
		height: 28px;
		align-items: center;
		justify-content: center;
		border: 0;
		border-radius: 6px;
		background: transparent;
		color: #64748b;
		cursor: pointer;
		font-family: inherit;
	}

	.meta-action:hover {
		background: #dbeafe;
		color: #2563eb;
	}

	.docs-article .typo {
		width: 100%;
	}

	.docs-article .ui.grid.m-margin-lr {
		margin-right: 0 !important;
		margin-left: 0 !important;
	}

	.docs-empty {
		padding: 30px 10px;
		color: #6b7280;
		text-align: center;
	}

	.docs-empty.error {
		color: #dc2626;
	}

	h1::before, h2::before, h3::before, h4::before, h5::before, h6::before {
		display: block;
		content: " ";
		height: 55px;
		margin-top: -55px;
		visibility: hidden;
	}

	@media (max-width: 768px) {
		.docs-page {
			display: block;
		}

		.docs-mobile-picker {
			display: block;
			margin-bottom: 12px;
		}

		.docs-mobile-picker .el-select {
			width: 100%;
		}

		.docs-content {
			width: 100%;
		}

		.docs-article {
			border-radius: 0 !important;
		}

		.blog-header {
			margin: 0 1rem 20px;
			padding: 0 0 20px;
		}

		.blog-header-top,
		.blog-header-left {
			align-items: flex-start;
			flex-direction: column;
		}

		.header-category {
			align-self: flex-start;
		}

		.blog-title {
			font-size: 28px;
		}
	}
</style>

<style>
	.docs-page .js-toc-content h1,
	.docs-page .js-toc-content h2 {
		display: flex;
		align-items: baseline;
		gap: 10px;
	}

	.docs-page .js-toc-content h3,
	.docs-page .js-toc-content h4,
	.docs-page .js-toc-content h5,
	.docs-page .js-toc-content h6 {
		padding-bottom: 0;
		border-bottom: 0;
	}

	.docs-page .docs-heading-number {
		flex: 0 0 auto;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-width: 32px;
		padding: 2px 8px;
		border-radius: 999px;
		background: #eff6ff;
		color: #2563eb;
		font-size: .72em;
		font-weight: 800;
		font-variant-numeric: tabular-nums;
	}
</style>

