import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { createPinia } from 'pinia';
import naive from 'naive-ui';

import App from './App.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/BookLibrary',
      name: 'BookLibrary',
      component: () => import('./pages/BookLibrary.vue'),
    },
    {
      path: '/About',
      name: 'About',
      component: () => import('./pages/About.vue'),
    },
    {
      path: '/WordLists',
      name: 'WordLists',
      component: () => import('./pages/WordLists.vue'),
    },
    {
      path: '/FlashCards',
      name: 'FlashCards',
      component: () => import('./pages/FlashCards.vue'),
    },
    {
      path: '/book/:bookID',
      name: 'BookStats',
      component: () => import('./pages/BookStats.vue'),
      props: true,
    },
  ],
});
console.log(router);

const pinia = createPinia();
const app = createApp(App);
app.use(router);
app.use(pinia);
// todo dont import everything?
app.use(naive);
app.mount('#app');
