<template>
  <div>
    <el-tabs v-model="activeName" @tab-change="changeTab">
      <el-tab-pane label="邮件" name="email"></el-tab-pane>
      <el-tab-pane label="Slack" name="slack"></el-tab-pane>
      <el-tab-pane label="Webhook" name="webhook"></el-tab-pane>
    </el-tabs>
    <pre>
      通知模板支持的变量：

        TaskId   任务ID
        TaskName 任务名称
        Status   任务执行结果状态
        Result   任务执行输出
    </pre>
  </div>
</template>

<script>
export default {
  data () {
    return {
      activeName: ''
    }
  },
  created () {
    const segments = this.$route.path.split('/')
    if (segments.length !== 4) {
      return 'email'
    }
    this.activeName = segments[3]
  },
  methods: {
    changeTab (name) {
      this.$router.push(`/system/notification/${name}`)
    }
  }
}
</script>
