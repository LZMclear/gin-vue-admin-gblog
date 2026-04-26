<template>
  <el-input
    v-model="value"
    type="textarea"
    :autosize="{ minRows: 10, maxRows: 24 }"
    resize="vertical"
    placeholder="请输入 Markdown 内容"
  />
</template>

<script setup>
  import { computed } from 'vue'

  const props = defineProps({
    modelValue: {
      type: String,
      default: ''
    }
  })

  const emit = defineEmits(['update:modelValue'])

  const value = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })
</script>
