<template>
  <div class="flex flex-col container mx-auto px-4 h-full">
    <ag-grid-vue
        class="ag-theme-alpine w-5/6 mx-auto
        h-full flex-grow-1 text-xl"
        :columnDefs="columnDefs"
        :rowData="rowData"
        @grid-ready="onGridReady"
        >
    </ag-grid-vue>
  </div>
</template>

<script setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { AgGridVue } from 'ag-grid-vue3';
import AddToCardQueue from '@/components/AddToCardQueue.vue';

const columnDefs = [
  {
    headerName: 'word',
    field: 'word',
    suppressSizeToFit: true,
  },
  {
    headerName: 'sentence',
    field: 'sentence',
    sort: 'desc',
  },
  {
    headerName: '',
    field: 'Make FlashCard',
    width: 50,
    cellRenderer: AddToCardQueue,
  },
];
let api = null;
// let columnApi = null;
function onGridReady(params) {
  api = params.api;
  // I know this will probably be used
  // columnApi = params.columnApi;
  api.sizeColumnsToFit();
}

const rowData = await window.ipc.loadFlaggedCards();
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
