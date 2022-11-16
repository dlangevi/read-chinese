<template>
  <div
    class="h-7/8
           container
           m-4
           mx-auto
           flex
           h-full
           w-3/4
           flex-col
           gap-2
           px-4"
  >
    <n-input
      placeholder="Input here (using 汉字 for now)"
      @input="onUpdateSearchBox"
    />
    <ag-grid-vue
      class="ag-theme-alpine mx-auto h-full
        w-full grow text-xl"
      :get-row-id="getRowId"
      :column-defs="columnDefs"
      :row-data="rowData"
      @grid-ready="onGridReady"
    />
  </div>
</template>

<script lang="ts" setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { AgGridVue } from 'ag-grid-vue3';
import { ref, onUnmounted } from 'vue';
import { NInput } from 'naive-ui';
import AddToCardQueue from '@/components/AddToCardQueue.vue';
import MarkLearned from '@/components/MarkLearned.vue';
import type { GetRowIdParams, GridReadyEvent, ColDef } from 'ag-grid-community';

import { backend } from '@wailsjs/models';

import { GetPossibleWords } from '@wailsjs/backend/Dictionaries';

const rowData = ref<backend.UnknownWordEntry[]>([]);
let currentSearch:string = '';
const getRowId = (params:GetRowIdParams) => params.data.word;

const columnDefs:ColDef[] = [
  {
    headerName: 'word',
    field: 'word',
    suppressSizeToFit: true,
    cellClass: 'text-xl',
    sort: 'asc',
    comparator: (wordA:string, wordB:string) => {
      let base = wordA.localeCompare(wordB);
      if (wordA.startsWith(currentSearch)) {
        base -= 10;
      }
      if (wordB.startsWith(currentSearch)) {
        base += 10;
      }
      return base;
    },
  },
  {
    headerName: 'pinyin',
    field: 'pinyin',
    width: 100,
    cellClass: [
      'border-2',
      'text-opacity-0',
      'text-black',
      'hover:text-opacity-100',
    ],
  },
  {
    headerName: 'definition',
    field: 'definition',
    minWidth: 400,
  },
  {
    headerName: '',
    field: 'markButton',
    width: 120,
    cellRenderer: MarkLearned,
  },
  {
    headerName: '',
    field: 'Make FlashCard',
    width: 120,
    cellRenderer: AddToCardQueue,
    cellRendererParams: {
      text: 'Create FlashCard',
      create: true,
    },
  },
];

let resizeCallback: () => void;
function onGridReady(params:GridReadyEvent) {
  params.api.sizeColumnsToFit();
  resizeCallback = () => {
    setTimeout(() => {
      params.api.sizeColumnsToFit();
    });
  };
  window.addEventListener('resize', resizeCallback);
  params.api.sizeColumnsToFit();
}

onUnmounted(() => {
  window.removeEventListener('resize', resizeCallback);
});

async function onUpdateSearchBox(newSearch:string) {
  currentSearch = newSearch;
  if (newSearch.length === 0) {
    rowData.value = [];
  } else {
    rowData.value = await GetPossibleWords(newSearch);
  }
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
