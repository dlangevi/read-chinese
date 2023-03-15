<template>
  <nav class="navbar h-16 bg-accent px-8">
    <div class="flex items-center justify-between">
      <router-link to="/Welcome">
        <img
          class="h-auto w-32"
          src="../assets/logo_transparent.png"
          alt="Read Chinese Logo"
        >
      </router-link>
      <div class="ml-6 flex space-x-4">
        <router-link
          v-for="item in navigation"
          :key="item.name"
          :to="item.href"
          :class="[{'bg-accent-focus': item.isActive},
                   'px-3 py-2 rounded-md text-accent-content',
                   'hover:bg-accent-focus',
                   'text-md font-semibold']"
          :aria-current="{'page': item.isActive}"
        >
          {{ item.name }}
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router';
import { computed } from 'vue';
import { getUserSettings } from '@/lib/userSettings';
const route = useRoute();
const UserSettings = getUserSettings();

type listEntry = {
  name: string,
  href: string,
  isActive?: boolean,
};

const navigation = computed(() => [
  { name: 'Book Library', href: '/BookLibrary' },
  { name: 'Manage Current Flashcards', href: '/FlashCards' },
  { name: 'Create New Flashcards', href: '/LearnWords' },
  ...(UserSettings.meta.EnableExperimental
    ? [{ name: 'View Known Words', href: '/KnownWords' }]
    : []),
  { name: 'Stats', href: '/Stats' },
  { name: 'Settings', href: '/Settings' },
].map((item:listEntry) => {
  item.isActive = item.href === route.fullPath;
  return item;
}));
</script>
