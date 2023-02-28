<template>
  <div class="mx-auto bg-accent px-8">
    <div class="flex h-16 items-center justify-between">
      <router-link to="/Welcome">
        <img
          class="block h-32 w-auto"
          src="../assets/logo_transparent.png"
          alt="Read More"
        >
      </router-link>
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
              :class="[item.href === route.fullPath
                         ? 'bg-accent-focus '
                         : 'hover:bg-accent-focus',
                       'px-3 py-2 rounded-md text-accent-content',
                       'text-sm font-medium']"
              :aria-current="item.href === route.fullPath ? 'page' : undefined"
            >
              {{ item.name }}
            </router-link>
          </div>
        </div>
      </div>
      <label for="my-drawer" class="btn-ghost btn-square btn">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          class="inline-block h-6 w-6 stroke-current"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 6h16M4 12h16M4 18h16"
        /></svg>
      </label>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router';
import { computed } from 'vue';
import { getUserSettings } from '@/lib/userSettings';
const UserSettings = getUserSettings();

const route = useRoute();
const navigation = computed(() => [
  { name: 'BookLibrary', href: '/BookLibrary' },
  ...(UserSettings.meta.EnableExperimental
    ? [{ name: 'Manage FlashCards', href: '/FlashCards' },
      { name: 'Create FlashCards', href: '/MakeCards' }]
    : []),
  { name: 'Wordlists', href: '/WordLists' },
  { name: 'Stats', href: '/Stats' },
  { name: 'Settings', href: '/Settings' },
]);

</script>
