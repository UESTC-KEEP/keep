import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import dpcreate from '@/components/dpcreate'
import create from '@/components/create'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'create',
      component:create
    },
    {
      path: '/form',
      name: 'getForm',
      component:dpcreate
    }
  ]
})
