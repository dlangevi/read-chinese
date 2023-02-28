<template>
  <div
    id="app"
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
          <router-view
            :key="route.fullPath"
            class="h-[calc(100vh-4rem)] w-screen overflow-scroll "
          />
        </div>
      </loader-provider>
    </message-provider>
  </div>
</template>

<script lang="ts" setup>
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';
import MessageProvider from '@/components/MessageProvider.vue';
import LoaderProvider from '@/components/LoaderProvider.vue';
import TopNav from '@/components/TopNav.vue';
import CardCreator from '@/components/CardCreator.vue';
import { themeChange } from 'theme-change';
import { getUserSettings } from '@/lib/userSettings';

import './App.css';

const UserSettings = getUserSettings();

onMounted(() => {
  themeChange(false);
});

const route = useRoute();

</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
