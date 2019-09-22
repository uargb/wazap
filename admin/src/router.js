import Vue from 'vue'
import Router from 'vue-router'
import Login from './views/Login.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'login',
      component: Login
    },
    {
      path: '/general',
      name: 'general',
      component: () => import(/* webpackChunkName: "general" */ './views/General.vue')
    },
    {
      path: '/qa',
      name: 'qa',
      component: () => import(/* webpackChunkName: "qa" */ './views/QA.vue')
    }
  ]
})
