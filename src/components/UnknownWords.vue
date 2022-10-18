<template>
  <ag-grid-vue
    class="ag-theme-alpine text-xl"
    :getRowId="getRowId"
    :columnDefs="columnDefs"
    :rowData="rowData"
    :context="gridContext"
    @grid-ready="onGridReady" />

</template>

<script lang="ts" setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { onBeforeMount, ref } from 'vue';
import { AgGridVue } from 'ag-grid-vue3';
import MarkLearned from '@components/MarkLearned.vue';
import AddToCardQueue from '@components/AddToCardQueue.vue';
import { getUserSettings } from '@/UserSettings';
import type { GetRowIdParams, GridReadyEvent, ColDef } from 'ag-grid-community';

const UserSettings = getUserSettings();

const props = defineProps<{
  bookFilter?: number,
}>();

const getRowId = (params:GetRowIdParams) => params.data.word;

const gridContext = {
  bookId: props.bookFilter,
};

function onGridReady(params:GridReadyEvent) {
  params.api.sizeColumnsToFit();
  window.addEventListener('resize', () => {
    setTimeout(() => {
      params.api.sizeColumnsToFit();
    });
  });
  params.api.sizeColumnsToFit();
}

const columnDefs:ColDef[] = [
  {
    headerName: 'word',
    field: 'word',
    width: 80,
    minWidth: 80,
    cellClass: 'text-xl',
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
    headerName: 'occurance',
    field: 'occurance',
    sort: 'desc',
    width: 50,
    minWidth: 50,
  },
];

const showDefinitions = UserSettings.Dictionaries.ShowDefinitions.read();
if (showDefinitions) {
  columnDefs.push(
    {
      headerName: 'definition',
      field: 'definition',
      minWidth: 400,
    },
  );
}

columnDefs.push(
  {
    headerName: '',
    field: 'markButton',
    width: 120,
    cellRenderer: MarkLearned,
  },
);
columnDefs.push(
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
);

const rowData = ref<any[]>([]);
onBeforeMount(async () => {
  if (props.bookFilter !== undefined) {
    rowData.value = await window.nodeIpc.learningTarget([props.bookFilter]);
  } else {
    rowData.value = await window.nodeIpc.learningTarget();
  }
  console.log(rowData);
});
</script>
