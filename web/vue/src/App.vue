<script setup>
import navMenu from './components/navMenu.vue'
import sidebar from './components/sidebar.vue'
</script>
<template>
  <el-container>
    <el-header v-cloak>
      <navMenu></navMenu>
    </el-header>
    <el-container v-cloak>
      <sidebar></sidebar>
      <el-main>
        <router-view/>
      </el-main>
    </el-container>
    <el-backtop :right="100" :bottom="100" />
  </el-container>
</template>

<script>
import installService from './api/install'

export default {
  created () {
    installService.status((data) => {
      if (!data) {
        this.$router.push('/install')
      }
    })
  }
}
</script>
<style>
  body {
    margin: 0;
  }
  [v-cloak] {
    display: none !important;
  }
  .el-switch {
    --el-switch-on-color: var(--el-color-success);
    --el-switch-off-color: var(--el-color-danger);
  }
  .el-table {
    margin: 20px 0;
  }
  .el-select {
    width: 192px; /* same as .el-input */
  }
</style>
