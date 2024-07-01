<template>
  <el-aside width="150px">
    <el-menu
      :default-active="currentRoute"
      mode="vertical"
      background-color="#334157"
      text-color="#fff"
      active-text-color="#ffd04b"
      router>
      <el-menu-item index="/host">
        <el-icon><Share /></el-icon>
        <span>节点列表</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script>
export default {
  name: 'host-sidebar',
  data () {
    return {}
  },
  computed: {
    currentRoute () {
      return '/host'
    }
  }
}
</script>
