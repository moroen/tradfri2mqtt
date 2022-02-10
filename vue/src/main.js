import { createApp } from 'vue'
import App from './App.vue'
import store from './store'
import router from './router'
import { Quasar, Notify, Dialog } from 'quasar'




// Import icon libraries
import '@quasar/extras/material-icons/material-icons.css'

// Import Quasar css
import 'quasar/src/css/index.sass'

// Assumes your root component is App.vue
// and placed in same folder as main.js

const myApp = createApp(App)

myApp.use(router)
myApp.use(store)

myApp.use(Quasar, {
  plugins: { Notify, Dialog }, // import Quasar plugins and add here
  /*
  config: {
    brand: {
      // primary: '#e46262',
      // ... or all other brand colors
    },
    notify: {...}, // default set of options for Notify Quasar plugin
    loading: {...}, // default set of options for Loading Quasar plugin
    loadingBar: { ... }, // settings for LoadingBar Quasar plugin
    // ..and many more (check Installation card on each Quasar component/directive/plugin)
  }
  */
})

// Assumes you have a <div id="app"></div> in your index.html
myApp.mount('#app')