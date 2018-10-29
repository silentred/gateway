// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import 'admin-lte/bootstrap/css/bootstrap.min.css'
import 'admin-lte/dist/css/AdminLTE.css'
import 'admin-lte/dist/css/skins/skin-blue.min.css'

//import 'admin-lte/plugins/jQuery/jquery-2.2.3.min.js'
import 'admin-lte/bootstrap/js/bootstrap.min.js'
import 'admin-lte/dist/js/app.min.js'

import Vue from 'vue'
import App from './App'
import router from './router'


Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App }
})
