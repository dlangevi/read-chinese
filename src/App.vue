<template>
  <div id="app" class="bg-white max-h-screen max-w-screen">
    <TopNav class="w-screen h-16"/>
    <div class="w-screen h-[calc(100vh-4rem)] overflow-scroll">
      <router-view />
    </div>
  </div>
</template>

<script>

import TopNav from './components/TopNav.vue';

export default {
  name: 'App',
  components: {
    TopNav,
  },
  data() {
    return {
    };
  },
  methods: {
  },
  mounted() {
  },
};
</script>

<style src="./app.css"/>
<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
