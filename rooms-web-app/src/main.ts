import { createApp } from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import { loadFonts } from './plugins/webfontloader'
import {key, store} from "./stores/store";

loadFonts()

createApp(App)
  .use(vuetify)
  .use(store,key)
  .mount('#app')
