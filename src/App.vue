<template>
  <n-message-provider>
    <div id="app" class="bg-white max-h-screen max-w-screen">
      <TopNav class="w-screen h-16" />
      <CardCreator />
      <div class="w-screen h-[calc(100vh-4rem)] overflow-auto">
        <Suspense>
          <router-view :key="route.fullPath" />
        </Suspense>
      </div>
    </div>
  </n-message-provider>
</template>

<script setup>
import TopNav from '@/components/TopNav.vue';
import CardCreator from '@/components/CardCreator.vue';
import { useRoute } from 'vue-router';
import { NMessageProvider } from 'naive-ui';
// eslint-disable-next-line import/no-unresolved
import 'virtual:windi.css';

const route = useRoute();
</script>

<style src="./app.css" />
<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
