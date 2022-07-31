import {createApp} from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import {loadFonts} from './plugins/webfontloader'
import {key, store} from "./stores/store";
import {stringFormat} from "@/plugins/string-formatter";

loadFonts()

createApp(App)
    .use(vuetify)
    .use(stringFormat)
    .use(store, key)
    .mount('#app')
