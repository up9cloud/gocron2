<template>
  <el-menu
    :default-active="currentRoute"
    mode="horizontal"
    :ellipsis="false"
    router>
    <el-menu-item index="/">
      <b style="font-size: var(--el-font-size-extra-large)">{{logo}}</b>
    </el-menu-item>
    <el-menu-item index="/task">
      <el-icon><Menu /></el-icon>
      <span>任务管理</span>
    </el-menu-item>
    <el-menu-item index="/host">
      <el-icon><Upload /></el-icon>
      <span>任务节点</span>
    </el-menu-item>
    <el-menu-item index="/user" v-if="$store.getters.user.isSuperAdmin">
      <el-icon><Service /></el-icon>
      <span>用户管理</span>
    </el-menu-item>
    <el-menu-item index="/system" v-if="$store.getters.user.isSuperAdmin">
      <el-icon><Setting /></el-icon>
      <span>系统管理</span>
    </el-menu-item>

    <el-menu-item @click="clickLogo" style='margin-left: auto;'>
      <el-icon v-if="!isDark"><Sunny /></el-icon>
      <el-icon v-if="isDark"><Moon /></el-icon>
    </el-menu-item>
    <el-sub-menu index="userStatus" v-if="$store.getters.user.token">
      <template #title>{{$store.getters.user.username}}</template>
      <el-menu-item index="/user/edit-my-password">
        <el-icon><Edit /></el-icon>
        <span>修改密码</span>
      </el-menu-item>
      <el-menu-item @click="logout" index="/user/logout">
        <el-icon><CircleClose /></el-icon>
        <span>退出</span>
      </el-menu-item>
    </el-sub-menu>
  </el-menu>
</template>

<script>
import { useDark, useToggle } from '@vueuse/core'
const isDark = useDark()
const toggleDark = useToggle(isDark)
export default {
  data () {
    return {
      logo: '🅖⚆🅒🅡🅞🅝②',
      isDark: isDark
    }
  },
  computed: {
    currentRoute () {
      if (this.$route.path === '/') {
        return '/task'
      }
      const segments = this.$route.path.split('/')
      return `/${segments[1]}`
    }
  },
  methods: {
    logout () {
      this.$store.commit('logout')
      this.$router.push('/')
    },
    clickLogo() {
      toggleDark()
    }
  }
}
</script>
