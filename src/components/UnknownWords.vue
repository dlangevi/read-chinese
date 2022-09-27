<template>
  <ag-grid-vue
    class="ag-theme-alpine text-xl"
    :getRowId="getRowId"
    :columnDefs="columnDefs"
    :rowData="rowData"
    :context="gridContext"
    @grid-ready="onGridReady" />

</template>

<script setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { onBeforeMount, ref, inject } from 'vue';
import { AgGridVue } from 'ag-grid-vue3';
import MarkLearned from '@/components/MarkLearned.vue';
import AddToCardQueue from '@/components/AddToCardQueue.vue';

const UserSettings = inject('userSettings');

const props = defineProps({
  showDefinitions: {
    type: Boolean,
    default: false,
  },
  bookFilter: {
    type: String,
    default: '',
  },
});
const getRowId = (params) => params.data.word;

const gridContext = {
  bookId: props.bookFilter,
};

function onGridReady(params) {
  params.api.sizeColumnsToFit();
  window.addEventListener('resize', () => {
    setTimeout(() => {
      params.api.sizeColumnsToFit();
    });
  });
  params.api.sizeColumnsToFit();
}

const columnDefs = [
  {
    headerName: 'word',
    field: 'word',
    width: 50,
    minWidth: 50,
    cellClass: 'text-xl',
  },
  {
    headerName: 'pinyin',
    field: 'pinyin',
    width: 50,
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
    width: 150,
    minWidth: 50,
  },
];

const showDefinitions = UserSettings.Dictionaries.ShowDefinitions.read();
console.log('show defs', showDefinitions);
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

const rowData = ref([]);
onBeforeMount(async () => {
  console.log(props.bookFilter);
  if (props.bookFilter !== '') {
    rowData.value = await window.ipc.learningTarget([props.bookFilter]);
  } else {
    rowData.value = await window.ipc.learningTarget();
  }
  console.log(rowData);
});
</script>
