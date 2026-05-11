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
					<div class="docs-meta" v-if="doc.path">
						<span>{{ doc.path }}</span>
						<a v-if="doc.editUrl" :href="doc.editUrl" target="_blank" rel="noopener noreferrer">
							<i class="github icon"></i>在 GitHub 上编辑
						</a>
					</div>
					<header v-if="doc.path" class="docs-doc-header">
						<h1>{{ docFileTitle }}</h1>
						<div class="docs-doc-stats">
							<span>作者：Gvto</span>
							<span>字数：{{ docWordCount }}</span>
							<span>阅读：约 {{ docReadMinutes }} 分钟</span>
						</div>
					</header>
					<div class="typo js-toc-content match-braces rainbow-braces" v-viewer v-html="doc.content"></div>
				</template>
			</div>
		</main>

		<aside class="docs-toc m-mobile-hide">
			<div class="docs-panel docs-toc-panel">
				<div class="docs-panel-title">
					<i class="list ul icon"></i>
					<button
						v-if="collapsibleTocItems.length"
						type="button"
						class="docs-toc-collapse-all"
						@click="toggleAllToc"
					>
						<i :class="isAllTocCollapsed ? 'angle double down icon' : 'angle double up icon'"></i>
					</button>
					<span>此页内容</span>
				</div>
				<div v-if="tocItems.length === 0" class="docs-toc-empty">暂无目录</div>
				<ul v-else class="docs-toc-list">
					<li
						v-for="item in visibleTocItems"
						:key="item.id"
						:class="['level-' + item.level, {active: item.id === activeTocId, collapsed: isTocCollapsed(item)}]"
					>
						<div class="docs-toc-row">
							<button
								v-if="item.level === 1 && hasTocChildren(item)"
								type="button"
								class="docs-toc-toggle"
								@click.stop="toggleToc(item)"
							>
								<i :class="isTocCollapsed(item) ? 'caret right icon' : 'caret down icon'"></i>
							</button>
							<span v-else class="docs-toc-toggle-placeholder"></span>
							<a href="" @click.prevent="scrollToHeading(item.id)">
								<span class="toc-number">{{ item.number }}</span>
								<span>{{ item.text }}</span>
							</a>
						</div>
					</li>
				</ul>
			</div>
		</aside>
	</div>
</template>

<script>
	import DocTree from '@/components/docs/DocTree'
	import {getDocContent, getDocsTree} from '@/api/docs'
	import {isSuccess} from '@/util/gvaResponse'

	export default {
		name: 'Docs',
		components: {DocTree},
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
				tocItems: [],
				activeTocId: '',
				collapsedTocIds: {}
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
			visibleTocItems() {
				const result = []
				let collapsedParent = ''
				this.tocItems.forEach(item => {
					if (item.level === 1) {
						collapsedParent = this.collapsedTocIds[item.id] ? item.id : ''
						result.push(item)
						return
					}
					if (!collapsedParent) {
						result.push(item)
					}
				})
				return result
			},
			collapsibleTocItems() {
				return this.tocItems.filter(item => item.level === 1 && this.hasTocChildren(item))
			},
			isAllTocCollapsed() {
				return this.collapsibleTocItems.length > 0
					&& this.collapsibleTocItems.every(item => this.collapsedTocIds[item.id])
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
		mounted() {
			window.addEventListener('scroll', this.updateActiveToc, {passive: true})
		},
		beforeDestroy() {
			window.removeEventListener('scroll', this.updateActiveToc)
		},
		beforeRouteLeave(to, from, next) {
			this.tocItems = []
			next()
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
				this.loadingContent = true
				this.error = ''
				this.tocItems = []
				this.activeTocId = ''
				this.collapsedTocIds = {}
				getDocContent(path).then(res => {
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
						this.buildToc()
						this.updateActiveToc()
					})
				}).catch(() => {
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
			buildToc() {
				const content = this.$el.querySelector('.js-toc-content')
				if (!content) {
					this.tocItems = []
					return
				}
				const usedIds = new Set()
				const counters = [0, 0]
				this.tocItems = Array.from(content.querySelectorAll('h1,h2'))
					.map((heading, index) => {
						const text = heading.textContent.trim()
						if (!text) {
							return null
						}
						const level = Number(heading.tagName.slice(1))
						const number = this.nextHeadingNumber(counters, level)
						let id = heading.getAttribute('id') || this.slugify(text) || `doc-heading-${index + 1}`
						const baseId = id
						let count = 1
						while (usedIds.has(id)) {
							count += 1
							id = `${baseId}-${count}`
						}
						usedIds.add(id)
						heading.setAttribute('id', id)
						this.applyHeadingNumber(heading, number)
						return {
							id,
							text,
							level,
							number
						}
					})
					.filter(Boolean)
			},
			hasTocChildren(item) {
				const index = this.tocItems.findIndex(tocItem => tocItem.id === item.id)
				return index >= 0 && this.tocItems[index + 1] && this.tocItems[index + 1].level === 2
			},
			isTocCollapsed(item) {
				return Boolean(this.collapsedTocIds[item.id])
			},
			toggleToc(item) {
				if (this.collapsedTocIds[item.id]) {
					this.$delete(this.collapsedTocIds, item.id)
				} else {
					this.$set(this.collapsedTocIds, item.id, true)
				}
			},
			toggleAllToc() {
				if (this.isAllTocCollapsed) {
					this.collapsedTocIds = {}
					return
				}
				const next = {}
				this.collapsibleTocItems.forEach(item => {
					next[item.id] = true
				})
				this.collapsedTocIds = next
			},
			nextHeadingNumber(counters, level) {
				const index = Math.max(0, Math.min(level - 1, counters.length - 1))
				counters[index] += 1
				for (let i = index + 1; i < counters.length; i += 1) {
					counters[i] = 0
				}
				return counters.slice(0, index + 1).filter(Boolean).join('.')
			},
			applyHeadingNumber(heading, number) {
				const oldNumber = heading.querySelector(':scope > .docs-heading-number')
				if (oldNumber) {
					oldNumber.remove()
				}
				const numberEl = document.createElement('span')
				numberEl.className = 'docs-heading-number'
				numberEl.textContent = number
				heading.insertBefore(numberEl, heading.firstChild)
			},
			scrollToHeading(id) {
				const heading = this.$el.querySelector(`#${window.CSS && CSS.escape ? CSS.escape(id) : id}`)
				if (!heading) {
					return
				}
				const top = heading.getBoundingClientRect().top + window.pageYOffset - 62
				window.scrollTo({
					top,
					behavior: 'smooth'
				})
				this.activeTocId = id
				if (this.$route.hash !== `#${id}`) {
					this.$router.replace({name: 'docs', query: this.$route.query, hash: `#${id}`}).catch(() => {})
				}
			},
			updateActiveToc() {
				const items = this.visibleTocItems
				if (!items.length) {
					return
				}
				let current = items[0].id
				for (const item of items) {
					const heading = this.$el.querySelector(`#${window.CSS && CSS.escape ? CSS.escape(item.id) : item.id}`)
					if (!heading) {
						continue
					}
					if (heading.getBoundingClientRect().top <= 90) {
						current = item.id
					} else {
						break
					}
				}
				this.activeTocId = current
			},
			slugify(text) {
				return text
					.trim()
					.toLowerCase()
					.replace(/\s+/g, '-')
					.replace(/[^\w\u4e00-\u9fa5-]/g, '')
			}
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

	.docs-sidebar,
	.docs-toc {
		position: sticky;
		top: 62px;
		max-height: calc(100vh - 82px);
		overflow: auto;
		scrollbar-width: thin;
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

	.docs-toc-collapse-all {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 26px;
		height: 26px;
		margin-left: auto;
		padding: 0;
		border: 1px solid #e5e7eb;
		border-radius: 4px;
		background: #fff;
		color: #64748b;
		cursor: pointer;
	}

	.docs-toc-collapse-all:hover {
		border-color: #bfdbfe;
		background: #eff6ff;
		color: #2563eb;
	}

	.docs-toc-collapse-all i {
		margin: 0 !important;
	}

	.docs-content {
		min-width: 0;
	}

	.docs-article {
		min-height: 520px;
		border: 1px solid #e5e7eb !important;
		border-radius: 6px !important;
		box-shadow: 0 10px 30px rgba(15, 23, 42, .035);
	}

	.docs-meta {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		padding-bottom: 14px;
		margin-bottom: 18px;
		border-bottom: 1px solid #e5e7eb;
		color: #6b7280;
		font-size: 13px;
	}

	.docs-meta span {
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.docs-meta a {
		color: #2563eb;
		white-space: nowrap;
	}

	.docs-doc-header {
		padding: 12px 0 24px;
		margin-bottom: 22px;
		border-bottom: 1px solid #edf2f7;
		text-align: center;
	}

	.docs-doc-header h1 {
		margin: 0 0 12px;
		color: #111827;
		font-size: 30px;
		line-height: 1.25;
		font-weight: 800;
	}

	.docs-doc-stats {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-wrap: wrap;
		gap: 8px 16px;
		color: #6b7280;
		font-size: 13px;
	}

	.docs-empty {
		padding: 30px 10px;
		color: #6b7280;
		text-align: center;
	}

	.docs-empty.error {
		color: #dc2626;
	}

	.docs-toc-empty {
		padding: 12px 4px;
		color: #9ca3af;
		font-size: 13px;
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

		.docs-meta {
			align-items: flex-start;
			flex-direction: column;
		}
	}
</style>

<style>
	.docs-toc-list {
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.docs-toc-list li {
		margin: 0;
	}

	.docs-toc-row {
		display: flex;
		align-items: flex-start;
		gap: 4px;
	}

	.docs-toc-toggle,
	.docs-toc-toggle-placeholder {
		flex: 0 0 18px;
		width: 18px;
		height: 30px;
	}

	.docs-toc-toggle {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0;
		border: 0;
		background: transparent;
		color: #64748b;
		cursor: pointer;
	}

	.docs-toc-toggle:hover {
		color: #2563eb;
	}

	.docs-toc-toggle i {
		margin: 0 !important;
	}

	.docs-toc-list a {
		flex: 1 1 auto;
		min-width: 0;
		display: flex;
		gap: 7px;
		padding: 6px 8px;
		border-left: 2px solid transparent;
		border-radius: 4px;
		color: #4b5563;
		font-weight: 300;
		line-height: 1.35;
		word-break: break-word;
	}

	.docs-toc-list .toc-number {
		flex: 0 0 auto;
		min-width: 22px;
		color: #2563eb;
		font-weight: 700;
		font-variant-numeric: tabular-nums;
	}

	.docs-toc-list a:hover,
	.docs-toc-list li.active a {
		background: #eff6ff;
		border-left-color: #2563eb;
		color: #2563eb !important;
		font-weight: 700;
	}

	.docs-toc-list .level-2 {
		padding-left: 10px;
	}

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
