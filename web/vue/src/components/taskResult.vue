<template>
  <el-tag :type="type" :effect="effect">{{ text }}</el-tag>
  <!-- <el-result :icon="type"></el-result> -->
</template>

<script>
const defaultValue = {
    type: 'info',
    text: '-',
    effect: 'plain',
}
function parse(row) {
    if (row.id === 0) {
        return defaultValue
    }
    switch (row.status) {
        case 0: return {
            type: 'danger',
            text: '失败',
            effect: "light",
        }
        case 1: return {
            type: 'success',
            text: '执行中',
            effect: 'plain',
        }
        case 2: return {
            type: 'success',
            text: '成功',
            effect: "light",
        }
        case 3: return {
            type: 'info',
            text: '取消',
            effect: "light",
        }
        default: return defaultValue
    }
}
export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return defaultValue
  },
  mounted() {
    this.parse()
  },
  updated() {
    this.parse()
  },
  methods: {
    parse() {
        const o = parse(this.modelValue)
        this.type = o.type
        this.text = o.text
        if (o.effect) {
            this.effect = o.effect
        }
    }
  }
}
</script>
<style scoped>
.el-result {
  --el-result-padding: 0;
  --el-result-icon-font-size: 75%;
}
</style>