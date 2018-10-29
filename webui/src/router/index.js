import Vue from 'vue'
import Router from 'vue-router'
import DashView from '@/components/Dash'
import ServiceList from '@/components/service/List'
import RouteList from '@/components/route/List'
import RouteForm from '@/components/route/Form'

Vue.use(Router)

export default new Router({
  //mode: 'history',
  scrollBehavior: function (to, from, savedPosition) {
    return savedPosition || { x: 0, y: 0 }
  },
  routes: [
    {
      path: '/',
      name: 'DashView',
      component: DashView,
      children: [
        {
          path: 'routes',
          alias: '',
          component: RouteList,
          name: 'RouteList',
          meta: {description: 'Route List'}
        },
        {
          path: 'route/create',
          component: RouteForm,
          name: 'RouteForm',
          meta: {description: 'Route Form'}
        },
        {
          path: 'services',
          component: ServiceList,
          name: 'ServiceList',
          meta: {description: 'Service List'}
        }
      ]
    }
  ]
})
