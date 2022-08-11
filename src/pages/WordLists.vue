<template>
  <div class="container mx-auto px-4">
    <div class="text-center">
      <h1 class="text-center text-xl mt-5">Maybe you know these words?</h1>
      <p>For now lets just mark learned words you already know</p>
    </div>
      <ag-grid-vue
          class="ag-theme-alpine w-full h-screen"
          :columnDefs="columnDefs"
          :rowData="rowData"
      >
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
        { headerName: 'word', field: 'word' },
        { headerName: 'occurance', field: 'occurance' },
        { headerName: 'markLearned', field: 'markButton', cellRenderer: MarkLearned },
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
