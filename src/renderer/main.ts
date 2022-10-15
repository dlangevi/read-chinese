import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { createPinia } from 'pinia';
import { generateUserSettings } from '@/shared/userSettings';
import { UserSettingsKey } from '@/shared/types';

import BookLibrary from '@/renderer/pages/BookLibrary.vue';
import AboutPage from '@/renderer/pages/AboutPage.vue';
import GlobalSettings from '@/renderer/pages/GlobalSettings.vue';
import UserStats from '@/renderer/pages/UserStats.vue';
import WordLists from '@/renderer/pages/WordLists.vue';
import FlashCards from '@/renderer/pages/FlashCards.vue';
import BookStats from '@/renderer/pages/BookStats.vue';

import App from './App.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/BookLibrary',
    },
    {
      path: '/BookLibrary',
      name: 'BookLibrary',
      component: BookLibrary,
    },
    {
      path: '/About',
      name: 'About',
      component: AboutPage,
    },
    {
      path: '/Settings',
      name: 'Settings',
      component: GlobalSettings,
    },
    {
      path: '/Stats',
      name: 'Stats',
      component: UserStats,
    },
    {
      path: '/WordLists',
      name: 'WordLists',
      component: WordLists,
    },
    {
      path: '/FlashCards',
      name: 'FlashCards',
      component: FlashCards,
    },
    {
      path: '/book/:bookId',
      name: 'BookStats',
      component: BookStats,
      props: true,
    },
  ],
});

async function init() {
  const pinia = createPinia();
  const app = createApp(App);
  const userSettings = await generateUserSettings();
  app.provide(UserSettingsKey, userSettings);
  app.use(router);
  app.use(pinia);
  app.mount('#app');
}
init();
