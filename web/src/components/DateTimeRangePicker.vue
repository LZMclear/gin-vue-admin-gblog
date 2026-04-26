<template>
  <el-date-picker
    v-model="localDate"
    type="datetimerange"
    range-separator="至"
    start-placeholder="开始时间"
    end-placeholder="结束时间"
    value-format="YYYY-MM-DD HH:mm:ss"
    @change="handleChange"
  />
</template>

<script setup>
  import { ref, watch } from 'vue'

  const props = defineProps({
    date: {
      type: Array,
      default: () => []
    },
    setDate: {
      type: Function,
      default: () => {}
    }
  })

  const localDate = ref(props.date)

  watch(
    () => props.date,
    (val) => {
      localDate.value = val
    }
  )

  const handleChange = (val) => {
    props.setDate(val || [])
  }
</script>
