<template>
  <div class="flex flex-col container mx-auto px-4 h-full">
    <div class="text-center flex-shrink-0">
      <h1 class="text-center text-xl mt-5">Maybe you know these words?</h1>
      <p>For now lets just mark learned words you already know 好不好?</p>
    </div>
    <unknown-words
      class="w-5/6 mx-auto h-full flex-grow-1"
      showDefinitions
      />
  </div>
</template>

<script setup>
import UnknownWords from '@/components/UnknownWords.vue';
</script>

<style scoped>
.ag-theme-alpine {
    --ag-foreground-color: rgb(126, 46, 132);
    --ag-background-color: rgb(249, 245, 227);
    --ag-header-foreground-color: rgb(204, 245, 172);
    --ag-header-background-color: rgb(209, 64, 129);
    --ag-odd-row-background-color: rgb(0, 0, 0, 0.03);
    --ag-header-column-resize-handle-color: rgb(126, 46, 132);
    --ag-font-size: 25px;
}
</style>
