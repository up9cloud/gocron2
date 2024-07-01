// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import { createApp } from 'vue'
import ElementPlus, { ElMessageBox } from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import store from './store/index'
import router from './router/index'

const app = createApp(App)
app.directive('focus', {
  inserted: function (el) {
    // 聚焦元素
    el.focus()
  }
})
app.config.globalProperties.$filters = {
  formatTime(time) {
    const fillZero = function (num) {
      return num >= 10 ? num : '0' + num
    }
    const date = new Date(time)
  
    const result = date.getFullYear() + '-' +
    (fillZero(date.getMonth() + 1)) + '-' +
    fillZero(date.getDate()) + ' ' +
    fillZero(date.getHours()) + ':' +
    fillZero(date.getMinutes()) + ':' +
    fillZero(date.getSeconds())
  
    if (result.indexOf('20') !== 0) {
      return ''
    }
  
    return result
  }
}
app.config.globalProperties.$appConfirm = function (callback) {
  this.$confirm('确定执行此操作?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    callback()
  })
}
app.use(store)
app.use(router)
app.use(ElementPlus)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.mount('#app')
