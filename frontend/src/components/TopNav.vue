<template>
  <nav class="navbar bg-accent px-8">
    <div class="flex h-16 items-center justify-between">
      <router-link to="/Welcome">
        <img
          class="block h-auto w-32"
          src="../assets/logo_transparent.png"
          alt="Read More"
        >
      </router-link>
      <div class="flex flex-1 items-center justify-center">
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
                       'text-md font-medium']"
              :aria-current="item.href === route.fullPath ? 'page' : undefined"
            >
              {{ item.name }}
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </nav>
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
