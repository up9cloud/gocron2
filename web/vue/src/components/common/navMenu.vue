<template>
  <div v-cloak>
    <el-menu
      :default-active="currentRoute"
      mode="horizontal"
      background-color="#334157"
      text-color="#fff"
      active-text-color="#ffd04b"
      router>
      <el-menu-item index="/" style="font-size:16px;letter-spacing:2px;color:rgb(255, 208, 75)">
        <b @click="changeLogo()">{{logos[logoIndex]}}</b>
      </el-menu-item>
      <el-menu-item index="/task">
        <el-icon><Menu /></el-icon>
        <span>任务管理</span>
      </el-menu-item>
      <el-menu-item index="/host">
        <el-icon><Upload /></el-icon>
        <span>任务节点</span>
      </el-menu-item>
      <el-menu-item index="/user" v-if="this.$store.getters.user.isSuperAdmin">
        <el-icon><Service /></el-icon>
        <span>用户管理</span>
      </el-menu-item>
      <el-menu-item index="/system" v-if="this.$store.getters.user.isSuperAdmin">
        <el-icon><Setting /></el-icon>
        <span>系统管理</span>
      </el-menu-item>

      <el-sub-menu index="userStatus" v-if="this.$store.getters.user.token" style='margin-left: auto;'>
        <template #title>{{this.$store.getters.user.username}}</template>
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
  </div>
</template>

<script>
export default {
  name: 'app-nav-menu',
  data () {
    return {
      logos: ['🅖⚆🅒🅡🅞🅝', 'Ⓖ🅞Ⓒ🅡Ⓞ🅝', '🅖➲🅒🅡🅞🅝', '🅶OCRON'],
      logoIndex: 0
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
    changeLogo () {
      this.logoIndex = this.logoIndex === (this.logos.length - 1) ? 0 : this.logoIndex + 1
    }
  }
}
</script>
