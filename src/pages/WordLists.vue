<style src="../app.css"/>
<template>
  <div class="flex flex-col container mx-auto px-4 h-full">
    <div class="text-center flex-shrink-0">
      <h1 class="text-center text-xl mt-5">Maybe you know these words?</h1>
      <p>For now lets just mark learned words you already know 好不好?</p>
    </div>
    <ag-grid-vue
        class="ag-theme-alpine w-1/2 mx-auto
        h-full flex-grow-1 text-xl"
        :columnDefs="columnDefs"
        :rowData="rowData" >
    </ag-grid-vue>
  </div>
</template>

<script>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { AgGridVue } from 'ag-grid-vue3';
import MarkLearned from '../components/MarkLearned.vue';

export default {
  name: 'WordLists',
  components: {
    AgGridVue,
  },
  data() {
    return {
      rowData: [
      ],
    };
  },
  setup() {
    return {
      columnDefs: [
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
      ],
    };
  },
  methods: {
  },

  async mounted() {
    const words = await window.ipc.learningTarget();
    console.log(words);
    this.rowData = words;
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
