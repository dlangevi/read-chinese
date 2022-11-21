<template>
  <div id="app" class="text-base-content">
    <n-message-provider>
      <CardCreator />
      <div class="flex h-screen w-screen flex-col bg-base-100">
        <TopNav class="w-screen basis-16" />
        <div class="w-screen grow basis-auto overflow-scroll">
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
