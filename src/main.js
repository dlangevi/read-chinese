import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import BookLibrary from './pages/BookLibrary.vue';
import BookStats from './pages/BookStats.vue';
import About from './pages/About.vue';

import App from './App.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'About',
      component: About,
    },
    {
      path: '/BookLibrary',
      name: 'BookLibrary',
      component: BookLibrary,
    },
    {
      path: '/About',
      name: 'About',
      component: About,
    },
    {
      path: '/book/:bookID',
      name: 'BookStats',
      component: BookStats,
      props: true,
    },
  ],
});
console.log(router);

const app = createApp(App);
app.use(router);
app.mount('#app');
