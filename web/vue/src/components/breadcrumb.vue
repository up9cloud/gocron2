<template>
  <el-breadcrumb separator-icon="ArrowRight" style="margin-bottom:20px">
    <el-breadcrumb-item v-for="l in layers">{{$te('routes.' + l.name) ? $t('routes.' + l.name) : l.id}}</el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script>
export default {
  computed: {
    layers() {
      const p = this.$route.path.split('/')
      let layers = [{
        path: '/',
        name: '--'
      }]
      for (let i = 2; i <= p.length; i++) {
        const s = p.slice(1, i).join('/')
        layers.push({
          path: '/' + s,
          name: s.replaceAll('/', '--'),
          id: p[i - 1]
        })
      }
      return layers
    },
  },
}
</script>