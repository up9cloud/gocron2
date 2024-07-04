import axios from 'axios'
import { ElMessage, ElLoading } from 'element-plus'
import router from '../router/index'
import store from '../store/index'
import Qs from 'qs'

// 成功状态码
const SUCCESS_CODE = 0
// 认证失败
const AUTH_ERROR_CODE = 401
// 应用未安装
const APP_NOT_INSTALL_CODE = 801

axios.defaults.baseURL = 'api'
axios.defaults.timeout = 10000
axios.defaults.responseType = 'json'

// let loading
// let count = 0
axios.interceptors.request.use(config => {
  config.headers['Auth-Token'] = store.getters.user.token
  // count++
  // if (!loading) {
  //   loading = ElLoading.service({
  //     lock: true,
  //   })
  // }
  return config
}, error => {
  return Promise.reject(error)
})
axios.interceptors.response.use(data => {
  return data
}, error => {
  return Promise.reject(error)
})

function handle (promise, next) {
  promise
    .then((res) => successCallback(res, next))
    .catch((error) => failureCallback(error))
    // .finally(() => {
    //   count--
    //   if (count === 0 && loading) {
    //     loading.close()
    //     loading = null
    //   }
    // })
}

function checkResponseCode (code, msg) {
  switch (code) {
    // 应用未安装
    case APP_NOT_INSTALL_CODE:
      router.push('/install')
      return false
    // 认证失败
    case AUTH_ERROR_CODE:
      router.push('/user/login')
      return false
    case SUCCESS_CODE:
      return true
    default:
      ElMessage({ message: msg })
      return false
  }
}

function successCallback (res, next) {
  if (!checkResponseCode(res.data.code, res.data.message)) {
    return
  }
  if (!next) {
    return
  }
  next(res.data.data, res.data.code, res.data.message)
}

function failureCallback (error) {
  ElMessage({
    showClose: true,
    message: '请求失败 - ' + error,
    type: 'error'
  })
}

export default {
  get (uri, params, next) {
    const promise = axios.get(uri, {params})
    handle(promise, next)
  },

  batchGet (uriGroup, next) {
    const requests = []
    for (let item of uriGroup) {
      let params = {}
      if (item.params !== undefined) {
        params = item.params
      }
      requests.push(axios.get(item.uri, {params}))
    }

    axios.all(requests)
      .then(axios.spread(function (...res) {
        const result = []
        for (let item of res) {
          if (!checkResponseCode(item.data.code, item.data.message)) {
            return
          }
          result.push(item.data.data)
        }
        next(...result)
      }))
      .catch((error) => failureCallback(error))
      // .finally(() => {
      //   count = count - requests.length
      //   if (count === 0 && loading) {
      //     loading.close()
      //     loading = null
      //   }
      // })
  },

  post (uri, data, next) {
    const promise = axios.post(uri, Qs.stringify(data), {
      headers: {
        post: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }
    })
    handle(promise, next)
  }
}
