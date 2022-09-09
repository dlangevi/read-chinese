<template>
<ag-grid-vue
    class="ag-theme-alpine text-xl"
    :getRowId="getRowId"
    :columnDefs="columnDefs"
    :rowData="rowData" >
</ag-grid-vue>

</template>

<script setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { onBeforeMount, ref } from 'vue';
import { AgGridVue } from 'ag-grid-vue3';
import MarkLearned from '@/components/MarkLearned.vue';

const props = defineProps({
  showDefinitions: {
    type: Boolean,
    default: false,
  },
  bookFilter: {
    type: Array,
    default: () => ([]),
  },
});
const getRowId = (params) => params.data.word;

const columnDefs = [
  {
    headerName: 'word',
    field: 'word',
  },
  {
    headerName: 'occurance',
    field: 'occurance',
    sort: 'desc',
  },
  {
    headerName: 'markLearned',
    field: 'markButton',
    cellRenderer: MarkLearned,
  },
];

if (props.showDefinitions) {
  columnDefs.push(
    {

      headerName: 'definition',
      field: 'definition',
      width: 600,
    },
  );
}

const rowData = ref([]);
onBeforeMount(async () => {
  console.log(props.bookFilter);
  rowData.value = await window.ipc.learningTarget(props.bookFilter);
});
</script>
