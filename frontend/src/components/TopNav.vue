<template>
  <div class="mx-auto bg-gray-800 px-8">
    <div class="flex h-16 items-center justify-between">
      <img
        class="block h-32 w-auto"
        src="../assets/logo_transparent.png"
        alt="Read More"
      >
      <div
        class="flex flex-1 items-center justify-center
          sm:items-stretch sm:justify-start"
      >
        <div class="hidden sm:ml-6 sm:block">
          <div class="flex space-x-4">
            <router-link
              v-for="item in navigation"
              :key="item.name"
              :to="item.href"
              :class="[item.current
                         ? 'bg-gray-900 text-white'
                         : 'text-gray-300 hover:bg-gray-700 hover:text-white',
                       'px-3 py-2 rounded-md text-sm font-medium']"
              :aria-current="item.current ? 'page' : undefined"
            >
              {{ item.name }}
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router';
import type { RouteLocationNormalizedLoaded } from 'vue-router';
import { watch, ref } from 'vue';

const route = useRoute();
// TODO get the image working better
const navigation = ref([
  { name: 'BookLibrary', href: '/BookLibrary', current: false },
  { name: 'Manage FlashCards', href: '/FlashCards', current: false },
  { name: 'Create FlashCards', href: '/MakeCards', current: false },
  { name: 'Manage Wordlist', href: '/WordLists', current: false },
  { name: 'Stats', href: '/Stats', current: false },
  { name: 'Settings', href: '/Settings', current: false },
  { name: 'About', href: '/About', current: false },
]);
function updateNav(currentRoute:RouteLocationNormalizedLoaded) {
  navigation.value.forEach((item) => {
    item.current = item.href === currentRoute.fullPath;
  });
}

watch(route, (_, newRoute) => {
  updateNav(newRoute);
});

</script>
