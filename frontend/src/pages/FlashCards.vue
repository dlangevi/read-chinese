<template>
  <div class="flex flex-col container mx-auto px-4 h-full">
    <ag-grid-vue
      class="ag-theme-alpine w-5/6 mx-auto
        h-full flex-grow text-xl"
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
import { onBeforeMount, ref, onUnmounted } from 'vue';
import { useMessage } from 'naive-ui';
import AddToCardQueue from '@/components/AddToCardQueue.vue';
import type { GetRowIdParams, GridReadyEvent, ColDef } from 'ag-grid-community';
import { LoadFlaggedCards } from '@wailsjs/backend/AnkiInterface';
import { backend } from '@wailsjs/models';

const message = useMessage();
const getRowId = (params:GetRowIdParams) => params.data.word;

const columnDefs:ColDef[] = [
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

const rowData = ref<backend.FlaggedCard[]>([]);
onBeforeMount(async () => {
  // TODO is there some way we can have the message be sent from
  // background, so we can just error at the source and not have
  // to do error handleing in all of the different vue files?
  try {
    rowData.value = await LoadFlaggedCards();
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
