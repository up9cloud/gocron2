<template>
  <el-empty description="您访问的页面不存在">
    <el-button type="primary" @click="jump">确定</el-button>
  </el-empty>
</template>

<script>
export default {
  methods: {
    jump () {
      this.$router.push('/')
    }
  }
}
</script>
