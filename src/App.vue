<template>
  <div id="app" class="bg-white">
    <button>
      <router-link to="/BookLibrary">Go to Library</router-link>
    </button>
    <router-link to="/about">Go to About</router-link>
    <router-view />
  </div>
</template>

<script>

export default {
  name: 'App',
  components: {
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
