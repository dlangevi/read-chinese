<template>
  <div data-theme="light">
    <n-message-provider>
      <div
        id="app"
        class="bg-base-100 max-w-screen max-h-screen"
      >
        <TopNav class="h-16 w-screen" />
        <CardCreator />
        <div class="h-[calc(100vh-4rem)] w-screen overflow-scroll">
          <Suspense>
            <router-view :key="route.fullPath" />
          </Suspense>
        </div>
      </div>
    </n-message-provider>
  </div>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router';
import { NMessageProvider } from 'naive-ui';
import TopNav from '@/components/TopNav.vue';
import CardCreator from '@/components/CardCreator.vue';

import 'virtual:windi.css';

const route = useRoute();
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
