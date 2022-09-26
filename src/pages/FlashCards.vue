<template>
  <div class="flex flex-col container mx-auto px-4 h-full">
    <ag-grid-vue
        class="ag-theme-alpine w-5/6 mx-auto
        h-full flex-grow text-xl"
        :getRowId="getRowId"
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
import { onBeforeMount, ref } from 'vue';
import { useMessage } from 'naive-ui';
import AddToCardQueue from '@/components/AddToCardQueue.vue';

const message = useMessage();
const getRowId = (params) => params.data.word;

const columnDefs = [
  {
    headerName: 'word',
    field: 'word',
    suppressSizeToFit: true,
    cellClass: 'text-xl',
  },
  {
    headerName: 'sentence',
    field: 'sentence',
    sort: 'desc',
    cellClass: 'text-xl',
  },
  {
    headerName: '',
    field: 'Make FlashCard',
    width: 50,
    cellRenderer: AddToCardQueue,
    cellRendererParams: {
      text: 'Replace Sentence',
    },
  },
];
let api = null;
// let columnApi = null;
function onGridReady(params) {
  api = params.api;
  // I know this will probably be used
  // columnApi = params.columnApi;
  api.sizeColumnsToFit();
  window.addEventListener('resize', () => {
    setTimeout(() => {
      api.sizeColumnsToFit();
    });
  });
  api.sizeColumnsToFit();
}

const rowData = ref([]);
onBeforeMount(async () => {
  // TODO is there some way we can have the message be sent from
  // background, so we can just error at the source and not have
  // to do error handleing in all of the different vue files?
  try {
    rowData.value = await window.ipc.loadFlaggedCards();
  } catch (e) {
    message.error('Please open Anki');
  }
});
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
