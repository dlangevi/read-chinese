import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { createPinia } from 'pinia';
// import naive from 'naive-ui';

import App from './App.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/BookLibrary',
      name: 'BookLibrary',
      component: () => { return import('./pages/BookLibrary.vue'); },
    },
    {
      path: '/About',
      name: 'About',
      component: () => { return import('./pages/AboutPage.vue'); },
    },
    {
      path: '/WordLists',
      name: 'WordLists',
      component: () => { return import('./pages/WordLists.vue'); },
    },
    {
      path: '/FlashCards',
      name: 'FlashCards',
      component: () => { return import('./pages/FlashCards.vue'); },
    },
    {
      path: '/book/:bookId',
      name: 'BookStats',
      component: () => { return import('./pages/BookStats.vue'); },
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
// app.use(naive);
app.mount('#app');
