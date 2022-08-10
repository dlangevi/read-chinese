<template>
  <div id="app" class="bg-white max-h-screen overflow-hidden">
  <Disclosure as="nav" class="bg-gray-800" v-slot="{ open }">
    <div class="max-w-7xl mx-auto px-2 sm:px-6 lg:px-8">
      <div class="relative flex items-center justify-between h-16">
        <div class="flex-1 flex items-center justify-center sm:items-stretch sm:justify-start">
          <div class="flex-shrink-0 flex items-center">
            <img class="block h-8 w-auto"
            src="https://tailwindui.com/img/logos/workflow-mark-indigo-500.svg" alt="Workflow" />
          </div>
          <div class="hidden sm:block sm:ml-6">
            <div class="flex space-x-4">
              <router-link v-for="item in navigation" :key="item.name" :to="item.href"
                :class="[item.current ? 'bg-gray-900 text-white' :
                                        'text-gray-300 hover:bg-gray-700 hover:text-white',
                                        'px-3 py-2 rounded-md text-sm font-medium']"
                :aria-current="item.current ? 'page' : undefined">{{ item.name }}
                </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Disclosure>
  <div class="overflow-scroll">
    <router-view />
  </div>
  </div>
</template>

<script>

import {
  Disclosure, DisclosureButton, DisclosurePanel, Menu, MenuButton, MenuItem, MenuItems,
} from '@headlessui/vue';
import { BellIcon, MenuIcon, XIcon } from '@heroicons/vue/outline';
import RouterButton from './components/RouterButton.vue';

export default {
  name: 'App',
  components: {
    RouterButton,
    Disclosure,
    DisclosureButton,
    DisclosurePanel,
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    BellIcon,
    MenuIcon,
    XIcon,
  },
  data() {
    return {
      // TODO have the current somehow match the current page
      navigation: [
        { name: 'BookLibrary', href: '/BookLibrary', current: true },
        { name: 'Manage FlashCards', href: '/FlashCards', current: true },
        { name: 'Manage Wordlist', href: '/WordLists', current: true },
        { name: 'About', href: 'About', current: true },
      ],
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
