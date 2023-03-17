<template>
  <div
    id="app-root"
    class="overflow-hide text-base-content"
    :data-theme="UserSettings.meta.Theme"
  >
    <message-provider>
      <loader-provider>
        <CardCreator />
        <div class="h-screen w-screen bg-base-100">
          <TopNav
            :class="['w-screen h-16',
                     {hidden: route.fullPath === '/Welcome'}]"
          />
          <!-- While we use Suspense, we want to limit async components to
          lightweight fetch only due to the current bug described here
          https://github.com/vuejs/router/issues/1324 -->
          <Suspense>
            <router-view
              :key="route.fullPath"
              :class="['w-screen overflow-auto',
                       {'h-[calc(100vh-4rem)]': route.fullPath !== '/Welcome'},
                       {'h-100vh': route.fullPath === '/Welcome'}
              ]"
            />
          </Suspense>
        </div>
      </loader-provider>
    </message-provider>
  </div>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router';
import MessageProvider from '@/components/MessageProvider.vue';
import LoaderProvider from '@/components/LoaderProvider.vue';
import TopNav from '@/components/TopNav.vue';
import CardCreator from '@/components/CardCreator.vue';
import { getUserSettings } from '@/lib/userSettings';

import './App.css';

const UserSettings = getUserSettings();
const route = useRoute();
</script>

<style>
#app-root {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
