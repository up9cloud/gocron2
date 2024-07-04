<template>
  <el-dialog v-model="value">
    <div>
      <pre>{{data.command}}</pre>
    </div>
    <el-divider></el-divider>
    <div>
      <pre>{{data.result}}</pre>
    </div>
  </el-dialog>
</template>
<script>
export default {
  props: ['modelValue', 'data'],
  emits: ['update:modelValue'],
  computed: {
    value: {
      get() {
        return this.modelValue
      },
      set(value) {
        this.$emit('update:modelValue', value)
      }
    }
  }
}
</script>