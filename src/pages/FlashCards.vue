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

<script>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { AgGridVue } from 'ag-grid-vue3';
import { useCardQueue } from '@/stores/CardQueue';
import AddToCardQueue from '@/components/AddToCardQueue.vue';

export default {
  name: 'FlashCards',
  components: {
    AgGridVue,
    AddToCardQueue,
  },
  data() {
    return {
      rowData: [],
      activeWord: '',
    };
  },
  setup() {
    const store = useCardQueue();
    return {
      store,
      columnDefs: [
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
          headerName: 'addToQueue',
          field: 'Make FlashCard',
          cellRenderer: AddToCardQueue,
        },
      ],
    };
  },
  methods: {
    addWord() {
      console.log('click3ed!');
      this.store.addWord('hello');
    },
    onGridReady(params) {
      this.api = params.api;
      this.columnApi = params.columnApi;
      this.api.sizeColumnsToFit();
    },

  },
  async mounted() {
    const flagged = await window.ipc.loadFlaggedCards();
    this.rowData = flagged;
  },
};
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
