/* global localStorage */

import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    username: localStorage.username === undefined ? '' : localStorage.username,
    password: localStorage.password === undefined ? '' : localStorage.password
  },
  mutations: {
    authenticate (state, payload) {
      state.username = payload.username
      localStorage.username = payload.username
      state.password = payload.password
      localStorage.password = payload.password
    },
    logout (state) {
      state.username = ''
      localStorage.username = ''
      state.password = ''
      localStorage.password = ''
    }
  }
})
