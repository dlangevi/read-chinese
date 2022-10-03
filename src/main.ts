import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { createPinia } from 'pinia';
import UserSettings from './userSettings';

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
      component: () => import('./pages/AboutPage.vue'),
    },
    {
      path: '/Settings',
      name: 'Settings',
      component: () => import('./pages/GlobalSettings.vue'),
    },
    {
      path: '/Stats',
      name: 'Stats',
      component: () => import('./pages/UserStats.vue'),
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
      path: '/book/:bookId',
      name: 'BookStats',
      component: () => import('./pages/BookStats.vue'),
      props: true,
    },
  ],
});

const pinia = createPinia();
const app = createApp(App);
app.provide('userSettings', UserSettings);
app.use(router);
app.use(pinia);
app.mount('#app');
