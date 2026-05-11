<template>
	<ul class="doc-tree">
		<li v-for="(node, index) in nodes" :key="nodeKey(node, index)">
			<div
				class="doc-tree-item"
				:class="{active: node.path && node.path === activePath, folder: isDir(node), collapsed: isDir(node) && isCollapsed(node, index)}"
				@click="handleClick(node, index)"
			>
				<i
					v-if="isDir(node)"
					:class="isCollapsed(node, index) ? 'caret right icon' : 'caret down icon'"
				></i>
				<i :class="node.type === 'dir' ? 'folder outline icon' : 'file alternate outline icon'"></i>
				<span>{{ node.title }}</span>
			</div>
			<DocTree
				v-if="node.children && node.children.length && !isCollapsed(node, index)"
				:nodes="node.children"
				:active-path="activePath"
				:parent-key="nodeKey(node, index)"
				@select="$emit('select', $event)"
			/>
		</li>
	</ul>
</template>

<script>
	export default {
		name: 'DocTree',
		props: {
			nodes: {
				type: Array,
				default: () => []
			},
			activePath: {
				type: String,
				default: ''
			},
			parentKey: {
				type: String,
				default: ''
			}
		},
		data() {
			return {
				collapsedDirs: {}
			}
		},
		methods: {
			handleClick(node, index) {
				if (this.isDir(node)) {
					this.toggleDir(node, index)
					return
				}
				if (node.type === 'file' && node.path) {
					this.$emit('select', node)
				}
			},
			isDir(node) {
				return node.type === 'dir'
			},
			nodeKey(node, index) {
				const rawKey = node.path || node.title || index
				return `${this.parentKey}/${node.type}:${rawKey}:${index}`
			},
			isCollapsed(node, index) {
				return Boolean(this.collapsedDirs[this.nodeKey(node, index)])
			},
			toggleDir(node, index) {
				const key = this.nodeKey(node, index)
				if (this.collapsedDirs[key]) {
					this.$delete(this.collapsedDirs, key)
				} else {
					this.$set(this.collapsedDirs, key, true)
				}
			}
		}
	}
</script>

<style scoped>
	.doc-tree {
		list-style: none;
		margin: 0;
		padding-left: 10px;
	}

	.doc-tree > li {
		margin: 3px 0;
	}

	.doc-tree-item {
		display: flex;
		align-items: center;
		gap: 8px;
		min-height: 32px;
		padding: 6px 9px;
		border-radius: 4px;
		color: #4b5563;
		cursor: pointer;
		line-height: 1.35;
		transition: background-color .16s ease, color .16s ease, transform .16s ease;
	}

	.doc-tree-item.folder {
		color: #374151;
		font-weight: 600;
	}

	.doc-tree-item.folder:hover {
		background: #f8fafc;
		color: #1f2937;
	}

	.doc-tree-item:not(.folder):hover,
	.doc-tree-item.active {
		background: #eff6ff;
		color: #2563eb;
	}

	.doc-tree-item:not(.folder):hover {
		transform: translateX(2px);
	}

	.doc-tree-item span {
		min-width: 0;
		word-break: break-word;
	}

	.doc-tree-item i {
		flex: 0 0 auto;
		margin: 0 !important;
	}
</style>
