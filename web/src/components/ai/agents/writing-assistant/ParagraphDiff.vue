<template>
  <div class="paragraph-diff">
    <div class="diff-toolbar">
      <el-radio-group v-model="viewMode" size="small">
        <el-radio-button value="inline">内联</el-radio-button>
        <el-radio-button value="split">分栏</el-radio-button>
      </el-radio-group>
      <div class="diff-actions">
        <el-button size="small" @click="setAll(true)">全部采用</el-button>
        <el-button size="small" @click="setAll(false)">全部撤销</el-button>
      </div>
    </div>
    <div class="diff-summary">已采纳 {{ stats.adopted }}/{{ stats.total }} 处修改</div>

    <div class="diff-list">
      <div
        v-for="(block, index) in blocks"
        :key="index"
        class="diff-block"
        :class="`is-${block.type}`"
      >
        <div class="block-tag">
          {{ typeLabel(block.type) }}
        </div>

        <template v-if="block.type === 'equal'">
          <pre class="block-text">{{ block.original }}</pre>
        </template>

        <template v-else-if="block.type === 'modified' && viewMode === 'split'">
          <div class="split-pane">
            <pre class="block-text original" :class="{ dimmed: block.takeRevised }">{{ block.original }}</pre>
            <pre class="block-text revised" :class="{ dimmed: !block.takeRevised }">{{ block.revised }}</pre>
          </div>
        </template>

        <template v-else-if="viewMode === 'split' && (block.type === 'added' || block.type === 'removed')">
          <div class="split-pane">
            <pre class="block-text original" :class="{ dimmed: block.takeRevised }">{{ block.original || '（无）' }}</pre>
            <pre class="block-text revised" :class="{ dimmed: !block.takeRevised }">{{ block.revised || '（已删除）' }}</pre>
          </div>
        </template>

        <template v-else>
          <pre v-if="block.type !== 'added'" class="block-text original" :class="{ dimmed: block.takeRevised && block.type !== 'equal' }">{{ block.original }}</pre>
          <pre v-if="block.type !== 'removed'" class="block-text revised" :class="{ dimmed: !block.takeRevised && block.type !== 'equal' }">{{ block.revised }}</pre>
        </template>

        <div v-if="block.type !== 'equal'" class="block-actions">
          <el-button
            size="small"
            :type="block.takeRevised ? 'primary' : 'default'"
            @click="choose(index, true)"
          >
            {{ block.type === 'added' ? '采用' : '采用 AI 版' }}
          </el-button>
          <el-button
            size="small"
            :type="!block.takeRevised ? 'primary' : 'default'"
            @click="choose(index, false)"
          >
            {{ block.type === 'removed' ? '保留' : '保留原文' }}
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { computed, ref } from 'vue'

  const props = defineProps({
    blocks: {
      type: Array,
      required: true
    }
  })

  const emit = defineEmits(['change'])

  const viewMode = ref('inline')

  const stats = computed(() => {
    const selectable = props.blocks.filter(
      (b) => b.type === 'modified' || b.type === 'added' || b.type === 'removed'
    )
    return {
      total: selectable.length,
      adopted: selectable.filter((b) => b.takeRevised).length
    }
  })

  const typeLabel = (type) => {
    const labels = {
      equal: '未改动',
      modified: '修改',
      added: '新增',
      removed: '删除'
    }
    return labels[type] || type
  }

  const choose = (index, takeRevised) => {
    props.blocks[index].takeRevised = takeRevised
    emit('change')
  }

  const setAll = (takeRevised) => {
    for (const block of props.blocks) {
      if (block.type === 'modified' || block.type === 'added') {
        block.takeRevised = takeRevised
      } else if (block.type === 'removed') {
        block.takeRevised = takeRevised // true=删除该段，false=保留
      }
    }
    emit('change')
  }
</script>

<style scoped lang="scss">
.paragraph-diff {
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.diff-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.diff-actions {
  display: flex;
  gap: 8px;
}

.diff-summary {
  margin: 6px 0;
  color: #909399;
  font-size: 12px;
}

.diff-list {
  flex: 1;
  min-height: 0;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.diff-block {
  position: relative;
  padding: 10px 12px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #fff;
}

.diff-block.is-modified {
  border-color: #e6a23c;
}

.diff-block.is-added {
  border-color: #67c23a;
}

.diff-block.is-removed {
  border-color: #f56c6c;
}

.block-tag {
  display: inline-block;
  margin-bottom: 6px;
  padding: 1px 8px;
  border-radius: 4px;
  background: #f0f2f5;
  color: #909399;
  font-size: 12px;
}

.is-modified .block-tag {
  background: #fdf6ec;
  color: #e6a23c;
}

.is-added .block-tag {
  background: #f0f9eb;
  color: #67c23a;
}

.is-removed .block-tag {
  background: #fef0f0;
  color: #f56c6c;
}

.block-text {
  margin: 0;
  padding: 6px 8px;
  border-radius: 4px;
  background: #fafbfc;
  color: #24292f;
  font-family: inherit;
  font-size: 13px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
}

.block-text.original {
  background: #fef0f0;
}

.block-text.revised {
  background: #f0f9eb;
}

.block-text.dimmed {
  opacity: 0.45;
  text-decoration: line-through;
}

.block-text.dimmed:empty {
  display: none;
}

.split-pane {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.block-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}

@media (max-width: 768px) {
  .split-pane {
    grid-template-columns: 1fr;
  }
}
</style>
