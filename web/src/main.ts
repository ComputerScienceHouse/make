import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { loadUser } from './auth.ts'

import '../node_modules/csh-material-bootstrap/dist/css/csh-material-bootstrap.min.css'
import 'bootstrap'
import 'bootstrap-icons/font/bootstrap-icons.css'

//import 'csh-material-bootstrap/color-modes.js';

async function load() {
  await loadUser()
  const app = createApp(App)
  app.use(router)
  app.mount('#app')
}

load()
