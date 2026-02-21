import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'; // Import the router
import { createPinia } from 'pinia'; // Import createPinia

const pinia = createPinia(); // Create the Pinia instance

createApp(App).use(router).use(pinia).mount('#app')
