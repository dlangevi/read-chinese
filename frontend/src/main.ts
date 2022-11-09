import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { createPinia } from 'pinia';
import { generateUserSettings } from '@/lib/userSettings';
import { UserSettingsKey } from '@/lib/types';
//
import GlobalSettings from '@/pages/GlobalSettings.vue';
import BookLibrary from '@/pages/BookLibrary.vue';
import AboutPage from '@/pages/AboutPage.vue';
import UserStats from '@/pages/UserStats.vue';
import WordLists from '@/pages/WordLists.vue';
import FlashCards from '@/pages/FlashCards.vue';
import MakeCards from '@/pages/MakeCards.vue';
import BookStats from '@/pages/BookStats.vue';

import App from './App.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/About',
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
      path: '/MakeCards',
      name: 'MakeCards',
      component: MakeCards,
    },
    {
      path: '/book/:bookId',
      name: 'BookStats',
      component: BookStats,
      props: (route) => {
        const bookIdParam = typeof route.params.bookId === 'object'
          ? route.params.bookId[0] : route.params.bookId;
        const bookId = Number.parseInt(bookIdParam, 10);
        if (Number.isNaN(bookId)) {
          console.error('failed to parse number');
          return 0;
        }
        return { bookId };
      },
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
