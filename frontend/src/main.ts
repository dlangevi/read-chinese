import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { createPinia } from 'pinia';
import { generateUserSettings, UserSettingsKey } from '@/lib/userSettings';
import GlobalSettings from '@/views/GlobalSettings.vue';
import Welcome from '@/views/Welcome.vue';
import BookLibrary from '@/views/BookLibrary.vue';
import AboutPage from '@/views/AboutPage.vue';
import UserStats from '@/views/UserStats.vue';
import WordLists from '@/views/WordLists.vue';
import FlashCards from '@/views/FlashCards.vue';
import MakeCards from '@/views/MakeCards.vue';
import BookStats from '@/views/BookStats.vue';

import App from './App.vue';

import { HealthCheck } from '@wailsjs/backend/Backend';
import { MessageApi, MessageApiKey } from '@/lib/messages';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/Welcome',
    },
    {
      path: '/Welcome',
      name: 'Welcome',
      component: Welcome,
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
          ? route.params.bookId[0]
          : route.params.bookId;
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
  const messageApi = new MessageApi();

  router.beforeEach(async (to, _from) => {
    console.log(to.name);
    const passes = await HealthCheck();
    if (!passes && to.name !== 'Welcome') {
      messageApi.error('Fix your problems before you can play');
      return '/Welcome';
    }
  });

  app.provide(MessageApiKey, messageApi);
  app.provide(UserSettingsKey, userSettings);
  app.use(router);
  app.use(pinia);
  app.config.errorHandler = (err, instance, info) => {
    console.log(err, instance, info);
    messageApi.error(`${err}`);
  };
  app.mount('#app');
}
init();
