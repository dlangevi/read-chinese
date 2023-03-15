<template>
  <div class="flex flex-col gap-4 p-4">
    <div class="mx-auto flex w-5/6 flex-col">
      Here are your known words
    </div>
    <ag-grid-vue
      class="ag-theme-alpine mx-auto
      w-5/6 grow text-xl"
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
import { ref, onUnmounted, onBeforeMount } from 'vue';
import type {
  GetRowIdParams,
  GridReadyEvent, ColDef,
} from 'ag-grid-community';
import { GetWordsGrid } from '@wailsjs/backend/knownWords';
import { backend } from '@wailsjs/models';

const getRowId = (params:GetRowIdParams) => params.data.word;

// TODO we want more columns/filters that can identify other problems
// some examples
// 1) Not in dict
// 2) Has flash card but no progress in X days
// 3) Had flashcard but no longer?
const columnDefs:ColDef[] = [
  {
    headerName: 'Word',
    field: 'Word',
    suppressSizeToFit: true,
    cellClass: 'text-xl',
  },
  {
    headerName: 'Interval',
    field: 'Interval',
    suppressSizeToFit: true,
    cellClass: 'text-xl',
    // TODO if interval === 10000 say 'Marked Known'
  },
  {
    headerName: 'Learned On',
    field: 'LearnedOn',
    suppressSizeToFit: true,
    cellClass: 'text-xl',
  },
  {
    headerName: 'In Dicts',
    field: 'InDict',
    suppressSizeToFit: true,
    cellClass: 'text-xl',
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

onBeforeMount(() => {
  loadData();
});

onUnmounted(() => {
  window.removeEventListener('resize', resizeCallback);
});

const rowData = ref<backend.WordGridRow[]>([]);

async function loadData() {
  rowData.value = await GetWordsGrid();
}
</script>
