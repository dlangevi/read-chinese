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
import {
  watch, onBeforeMount, ref, toRaw, onUnmounted,
} from 'vue';
import { AgGridVue } from 'ag-grid-vue3';
import MarkLearned from '@/components/MarkLearned.vue';
import AddToCardQueue from '@/components/AddToCardQueue.vue';
import { getUserSettings } from '@/lib/userSettings';
import type { GetRowIdParams, GridReadyEvent, ColDef } from 'ag-grid-community';
import { backend } from '@wailsjs/models';
import { GetDefinitions } from '@wailsjs/backend/Dictionaries';

const UserSettings = getUserSettings();

const props = defineProps<{
  words: backend.UnknownWordEntry[],
  bookFilter?: number,
}>();

const getRowId = (params:GetRowIdParams) => params.data.word;

const gridContext = {
  bookId: props.bookFilter,
};

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
watch(() => props.words, async (newWords) => {
  rowData.value = await GetDefinitions(toRaw(newWords));
  console.log('new Words', rowData.value);
});

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

onBeforeMount(async () => {
  const rawWords = toRaw(props.words);
  rowData.value = await GetDefinitions(rawWords);
  console.log(rowData);
});

</script>
