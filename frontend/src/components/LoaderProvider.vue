<template>
  <div>
    <div
      v-if="api.shouldShow"
      class="bg-opacity/50 fixed inset-0 z-[1000]
  flex items-center justify-center bg-gray-800"
    >
      <div class="flex flex-col items-center">
        <div
          class="mb-4 h-12 w-12 animate-spin
      rounded-full border-b-2 border-white"
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
