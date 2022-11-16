<template>
  <div class="text-base-content">
    <n-message-provider>
      <div
        id="app"
        class="max-w-screen max-h-screen bg-base-100"
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
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { NMessageProvider } from 'naive-ui';
import TopNav from '@/components/TopNav.vue';
import CardCreator from '@/components/CardCreator.vue';

import './App.css';

onMounted(() => {
  const theme = localStorage.getItem('theme');
  if (theme) {
    document.documentElement.setAttribute('data-theme', theme);
  }
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
