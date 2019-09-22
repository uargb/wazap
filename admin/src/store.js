import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    username: '',
    password: ''
  },
  mutations: {
    authenticate (state, payload) {
      state.username = payload.username
      state.password = payload.password
    },
    logout (state) {
      state.username = ''
      state.password = ''
    }
  }
})
