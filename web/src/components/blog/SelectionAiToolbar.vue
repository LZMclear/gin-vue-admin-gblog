<template>
  <Teleport to="body">
    <div
      v-if="visible && enabled"
      ref="toolbar"
      class="selection-ai-toolbar"
      role="group"
      aria-label="选区 AI 操作"
      :style="{ left: `${position.left}px`, top: `${position.top}px` }"
      @pointerdown.prevent
      @focusout="handleBlur"
    >
      <button type="button" @click="run('polish')">AI 润色</button>
      <button type="button" @click="run('rewrite')">AI 改写</button>
    </div>
  </Teleport>
</template>

<script setup>
  import { onBeforeUnmount, ref, watch } from 'vue'
  import { placeSelectionToolbar, textareaSelectionPoint } from './selectionToolbarPosition.js'

  const props = defineProps({ textarea: { type: Object, default: null }, enabled: Boolean })
  const emit = defineEmits(['action'])
  const toolbar = ref(null)
  const visible = ref(false)
  const position = ref({ left: 0, top: 0 })
  let frame = null
  let selecting = false
  const hide = () => {
    visible.value = false
    if (frame !== null) cancelAnimationFrame(frame)
    frame = null
  }
  const update = () => {
    frame = null
    const textarea = props.textarea
    if (!props.enabled || selecting || !textarea || (document.activeElement !== textarea && !toolbar.value?.contains(document.activeElement)) ||
      textarea.selectionStart === textarea.selectionEnd || !textarea.value.slice(textarea.selectionStart, textarea.selectionEnd).trim()) return hide()
    const point = textareaSelectionPoint(textarea)
    if (!point || point.y < 0 || point.y > window.innerHeight) return hide()
    position.value = placeSelectionToolbar(point, { width: window.innerWidth, height: window.innerHeight })
    visible.value = true
  }
  const schedule = (event) => {
    if (event?.key === 'Escape') return
    if (frame === null) frame = requestAnimationFrame(update)
  }
  const handleBlur = (event) => {
    if (event.relatedTarget === props.textarea || toolbar.value?.contains(event.relatedTarget)) return
    hide()
  }
  const escape = (event) => { if (event.key === 'Escape') hide() }
  const startSelection = () => { selecting = true; hide() }
  const finishSelection = () => {
    if (!selecting) return
    selecting = false
    schedule()
  }
  const run = (action) => {
    if (!props.enabled) return
    // 同步捕获选区，之后再打开抽屉；不依赖抽屉获取焦点后的 DOM 状态。
    emit('action', action)
    hide()
  }
  watch(() => props.textarea, (textarea, _, onCleanup) => {
    if (!textarea) return
    const events = { select: schedule, keyup: schedule, pointerdown: startSelection, input: hide, blur: handleBlur }
    for (const [name, handler] of Object.entries(events)) textarea.addEventListener(name, handler)
    document.addEventListener('pointerup', finishSelection)
    document.addEventListener('pointercancel', finishSelection)
    document.addEventListener('selectionchange', schedule)
    document.addEventListener('keydown', escape)
    window.addEventListener('scroll', hide, true)
    window.addEventListener('resize', hide)
    onCleanup(() => {
      for (const [name, handler] of Object.entries(events)) textarea.removeEventListener(name, handler)
      document.removeEventListener('pointerup', finishSelection)
      document.removeEventListener('pointercancel', finishSelection)
      document.removeEventListener('selectionchange', schedule)
      document.removeEventListener('keydown', escape)
      window.removeEventListener('scroll', hide, true)
      window.removeEventListener('resize', hide)
      selecting = false
      hide()
    })
  }, { immediate: true })
  watch(() => props.enabled, enabled => { if (!enabled) hide() })
  onBeforeUnmount(hide)
</script>

<style scoped>
.selection-ai-toolbar {
  position: fixed;
  z-index: 3001;
  display: flex;
  width: 184px;
  height: 38px;
  padding: 3px;
  box-sizing: border-box;
  border: 1px solid var(--el-border-color-light, #ddd);
  border-radius: 7px;
  background: var(--el-bg-color-overlay, #fff);
  box-shadow: 0 4px 16px #0002;
}
.selection-ai-toolbar button {
  flex: 1;
  border: 0;
  border-radius: 4px;
  background: transparent;
  color: var(--el-color-primary, #409eff);
  font: inherit;
  font-size: 13px;
  cursor: pointer;
}
.selection-ai-toolbar button:hover { background: var(--el-color-primary-light-9, #ecf5ff); }
.selection-ai-toolbar button:focus-visible { outline: 2px solid var(--el-color-primary, #409eff); }
</style>
