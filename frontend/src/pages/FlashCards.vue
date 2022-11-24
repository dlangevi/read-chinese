<template>
  <div class="container m-4 mx-auto flex h-full flex-col px-4">
    <ag-grid-vue
      class="ag-theme-alpine mx-auto h-full
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
import { onBeforeMount, ref, onUnmounted } from 'vue';
import { useMessage } from '@/lib/messages';
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
