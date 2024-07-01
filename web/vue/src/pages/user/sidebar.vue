<template>
  <el-aside width="150px">
    <el-menu
      :default-active="currentRoute"
      mode="vertical"
      background-color="#334157"
      text-color="#fff"
      active-text-color="#ffd04b"
      router>
      <el-menu-item index="/user">
        <el-icon><Notification /></el-icon>
        <span>用户列表</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script>
export default {
  name: 'user-sidebar',
  data () {
    return {}
  },
  computed: {
    currentRoute () {
      return '/user'
    }
  }
}
</script>
