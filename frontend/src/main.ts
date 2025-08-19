import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import './router/main.css'; // Asegúrate de importar tu CSS de Tailwind aquí

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.mount('#app');