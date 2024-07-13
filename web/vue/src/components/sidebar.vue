<template>
  <el-aside width="150px">
    <el-menu
      :default-active="currentLevel2Route"
      mode="vertical"
      router>
      <el-menu-item v-for="r in currentLevel2Routes" :index="r.path">
        <!-- <el-icon><BellFilled /></el-icon> -->
        <span>{{$t( 'routes.' + r.name)}}</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script>
import { routes } from '../router/index.js'
export default {
  computed: {
    currentLevel2Route() {
      return this.$route.path.split('/').splice(0, 3).join('/')
    },
    currentLevel2Routes() {
      const p = this.$route.path.replace(/^\//g, '').split('/')
      const filtered = routes.filter(r => {
        if (r.meta && r.meta.inSidebar === false) {
          return false
        }
        const pp = r.path.replace(/^\//g, '').split('/')
        if (pp.length !== 2) {
          return false
        }
        return pp[0] === p[0]
      })
      return filtered
    }
  },

}
</script>
