<template>
  <el-container>
    <el-header height="56px">
      <app-header></app-header>
      <app-nav-menu></app-nav-menu>
    </el-header>
    <el-main >
      <div id="main-container" v-cloak>
        <router-view/>
      </div>
    </el-main>
  </el-container>
</template>

<script>
import installService from './api/install'
import appHeader from './components/common/header.vue'
import appNavMenu from './components/common/navMenu.vue'
import appFooter from './components/common/footer.vue'

export default {
  name: 'App',
  data () {
    return {}
  },
  created () {
    installService.status((data) => {
      if (!data) {
        this.$router.push('/install')
      }
    })
  },
  components: {
    appHeader,
    appNavMenu,
    appFooter
  }
}
</script>
<style>
  [v-cloak] {
    display: none !important;
  }
  body {
    margin:0;
  }
  .el-switch {
    --el-switch-on-color: #13ce66;
    --el-switch-off-color: #ff4949
  }
  .el-header {
    padding:0;
    margin:0;
    border-bottom: 1px solid #fff;
  }
  .el-container {
    padding:0;
    margin:0;
    width: 100%;
    background-color: #eff1f7;
  }
  .el-main {
    padding:0;
    margin:0;
  }
  #main-container .el-main {
    height: calc(100vh - 116px);
    margin: 20px 20px 40px 20px;
  }
  .el-aside .el-menu {
    height: 100%;
  }
  .el-table {
    box-shadow:0 4px 8px 0 rgb(0 0 0 / 5%);
    color: #385879;
  }
  .el-pagination {
    background-color: #fff;
    border-radius: 1.5rem;
    padding: 6px 20px;
    box-shadow: 0 2px 4px 0 rgb(0 0 0 / 5%);
  }
  .el-pagination.is-background .btn-next, .el-pagination.is-background .btn-prev, .el-pagination.is-background .el-pager li {
    border-radius: 50%;
  }
  .el-form.el-form--inline .el-input {
    box-shadow: 0 2px 4px 0 rgb(0 0 0 / 3%);
  }
  .el-pagination select,
  .el-pagination input,
  .el-pagination span,
  .el-form input,
  .el-form label {
    color: #385879;
  }
  .el-form.page-form {
    background-color: #fff;
    padding: 20px 20px;
    box-shadow: 0 2px 4px 0 rgb(0 0 0 / 5%);
  }
  .el-button {
    box-shadow:0 3px 1px -2px rgb(0 0 0 / 7%), 0 2px 2px 0 rgb(0 0 0 / 7%), 0 1px 5px 1px rgb(0 0 0 / 7%);
  }
  .el-button.box-shadow-not {
    box-shadow: 0 0 0 0 rgb(0 0 0 / 0%);
  }
  .el-table th {
    background-color: #f5f7fa;
    color: #687a8d;
  }
  .el-breadcrumb .is-link{
    color: #536170;
  }
</style>
