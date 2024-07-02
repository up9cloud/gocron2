<template>
  <el-tag :type="type" :effect="effect">{{ text }}</el-tag>
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
            effect: "dark",
        }
        case 1: return {
            type: 'success',
            text: '执行中',
            effect: 'light',
        }
        case 2: return {
            type: 'success',
            text: '成功',
            effect: "dark",
        }
        case 3: return {
            type: 'info',
            text: '取消',
            effect: "dark",
        }
        default: return defaultValue
    }
}
export default {
  name: 'task-result',
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