<template>
  <div>
    <div
      v-if="api.shouldShow"
      class="
      fixed inset-0
      z-50 flex
      items-center justify-center
      bg-gray-900/50 "
    >
      <div class="flex flex-col items-center">
        <div
          class="my-4 h-24 w-24 animate-spin
        rounded-full border-b-4 border-white"
        />
        <p class="text-white">{{ api.loadingText }}</p>
      </div>
    </div>
    <slot />
  </div>
</template>

<script lang="ts" setup>
import { useLoader } from '@/lib/loading';
const api = useLoader();
</script>
