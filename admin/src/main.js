import Vue from 'vue'

import Buefy from 'buefy'
import 'buefy/dist/buefy.css'
import 'bulma-helpers/css/bulma-helpers.min.css'

import qs from 'qs'

import axios from 'axios'

import App from './App.vue'
import router from './router'
import store from './store'
Vue.prototype.$qs = qs
Vue.prototype.$axios = axios.create({
  baseURL: 'http://localhost:8090/',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded'
  }
})

Vue.use(Buefy)
Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
