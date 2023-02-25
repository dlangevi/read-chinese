<template>
  <div id="app" class="text-base-content">
    <message-provider>
      <loader-provider>
        <CardCreator />
        <div class="flex h-screen w-screen flex-col bg-base-100">
          <TopNav
            :class="['w-screen basis-16',
                     {hidden: route.fullPath === '/Welcome'}]"
          />
          <div class="w-screen grow basis-auto overflow-auto">
            <Suspense>
              <router-view :key="route.fullPath" />
            </Suspense>
          </div>
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
import { getUserSettings } from '@/lib/userSettings';

import './App.css';

const UserSettings = getUserSettings();

onMounted(() => {
  document.documentElement.setAttribute('data-theme', UserSettings.meta.Theme);
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
