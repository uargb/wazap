import Vue from 'vue'

import Buefy from 'buefy'
import 'buefy/dist/buefy.css'
import 'bulma-helpers/css/bulma-helpers.min.css'

import qs from 'qs'

import App from './App.vue'
import router from './router'
import store from './store'
Vue.prototype.$qs = qs
Vue.prototype.$apiBase = (login, password, url) => { return 'http://3.122.55.78:8090/admin/' + login + '/' + password + '/' + url }
Vue.prototype.$error = (self, text, duration = 5000) => {
  self.$buefy.toast.open({
    message: text,
    type: 'is-danger',
    duration: duration
  })
}

Vue.use(Buefy)
Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
