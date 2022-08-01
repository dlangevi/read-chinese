<template>
  <div id="app" class="bg-white">
    <RouterButton to="/BookLibrary" text="Library" />
    <RouterButton to="/FlashCards" text="Manage FlashCards" />
    <RouterButton to="/WordLists" text="Manage Wordlist" />
    <RouterButton to="/About" text="About" />
    <router-view />
  </div>
</template>

<script>
import RouterButton from './components/RouterButton.vue';

export default {
  name: 'App',
  components: {
    RouterButton,
  },
  data() {
    return {};
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
  margin-top: 60px;
}
</style>
