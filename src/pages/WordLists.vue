<template>
  <div class="flex h-full">
    <div class="flex-shrink-0 m-4">
      <n-space vertical>
      <n-button type="primary" @click="importLegacy">
        Import Legacy Words
      </n-button>
      <n-button type="primary" @click="importCSV">
        Import CSV file
      </n-button>
      <n-button type="primary" @click="importAnki">
        Sync from Anki
      </n-button>
      </n-space>

    </div>
    <div class="flex flex-col w-full h-full">
      <div class="text-center flex-shrink-0">
        <h1 class="text-center text-xl mt-5">Maybe you know these words?</h1>
        <p>For now lets just mark learned words you already know 好不好?</p>
      </div>
      <unknown-words
        class="w-5/6 mx-auto h-full flex-grow-1"
        showDefinitions
        />
    </div>
  </div>
</template>

<script setup>
import { NButton, NSpace } from 'naive-ui';
import UnknownWords from '@/components/UnknownWords.vue';

function importLegacy() {
  window.ipc.importLegacyWords();
}

function importCSV() {
  window.ipc.importCSVWords();
}

function importAnki() {
  window.ipc.importAnkiKeywords();
}
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
