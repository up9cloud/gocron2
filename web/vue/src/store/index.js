import { createStore } from 'vuex'
import userStorage from '../storage/user'
import router from '../router/index'

export default createStore({
  state: {
    user: userStorage.get()
  },
  getters: {
    user (state) {
      return state.user
    },
    login (state) {
      return state.user.token !== ''
    }
  },
  mutations: {
    setUser (state, user) {
      userStorage.setToken(user.token)
      userStorage.setUid(user.uid)
      userStorage.setUsername(user.username)
      userStorage.setIsAdmin(user.isAdmin)
      state.user = userStorage.get()
    },
    logout (state) {
      userStorage.clear()
      state.user = userStorage.get()
      router.push('/user/login')
    }
  }
})
