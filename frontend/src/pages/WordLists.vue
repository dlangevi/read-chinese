<template>
  <div class="flex h-full">
    <div class="flex-shrink-0 m-4">
      <n-space vertical>
        <n-button type="primary" @click="importCSV">
          Import CSV file
        </n-button>
        <n-button type="primary" @click="importAnki">
          Sync from Anki
        </n-button>
        <n-cascader
          :options="options"
          placeholder="Load HSK Words"
          :show-path="true"
          expand-trigger="click"
          check-strategy="child"
          @update:value="loadHsk"
        />
      </n-space>
    </div>
    <div class="flex flex-col w-full h-full">
      <div class="text-center flex-shrink-0">
        <h1 class="text-center text-xl mt-5">Maybe you know these words?</h1>
        <p>For now lets just mark learned words you already know 好不好?</p>
      </div>
      <unknown-words
        class="w-5/6 mx-auto h-full flex-grow"
        showDefinitions
        :words="words"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { NButton, NSpace, NCascader } from 'naive-ui';
import type { CascaderOption } from 'naive-ui';
import { ref, onBeforeMount } from 'vue';
import type {
  UnknownWordEntry, HskLevel, HskVersion,
} from '@/shared/types';
import { LearningTarget } from '@wailsjs/backend/BookLibrary';
import UnknownWords from '../components/UnknownWords.vue';

interface HskCascaderOption extends CascaderOption {
  level: HskLevel,
  version: HskVersion,
}
const options:CascaderOption[] = ['2.0', '3.0'].map((version) => ({
  label: `HSK ${version}`,
  value: version,
  children: [1, 2, 3, 4, 5, 6].map((level) => ({
    value: `${version}-${level}`,
    version,
    level,
    label: `Level ${level}`,
  })),
}));

const words = ref<UnknownWordEntry[]>([]);
onBeforeMount(async () => {
  words.value = await LearningTarget();
});

function importCSV() {
  window.nodeIpc.importCSVWords();
}

function importAnki() {
  window.nodeIpc.importAnkiKeywords();
}

async function loadHsk(_: string, option: HskCascaderOption) {
  words.value = await window.nodeIpc.hskWords(option.version, option.level);
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
}
</style>
