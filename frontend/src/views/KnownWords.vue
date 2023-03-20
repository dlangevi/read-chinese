<template>
  <div class="flex gap-4 p-4">
    <div class="mx-auto flex w-1/5 flex-col gap-4">
      <h2 class="text-3xl">
        Here are your known words
      </h2>
      <div> Total words: {{ rowData.length }} </div>
      <div>
        Words not in dicts: {{
          rowData.filter(row => !row.InDict).length
        }}
      </div>
      <div>
        Words marked learned: {{
          rowData.filter(row => row.Interval === 10000).length
        }}
      </div>
      <button class="btn-primary btn" @click="importAnki">
        Sync from Anki
      </button>
      <button class="btn-primary btn" @click="importAnkiReviews">
        Sync dates from Anki
      </button>
      <import-csv v-if="UserSettings.meta.EnableExperimental" />
    </div>
    <ag-grid-vue
      class="ag-theme-alpine mx-auto
      w-3/5 text-xl"
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
import { ref, onUnmounted, onBeforeMount } from 'vue';
import type {
  GetRowIdParams, ValueGetterParams,
  GridReadyEvent, ColDef,
} from 'ag-grid-community';
import { GetWordsGrid } from '@wailsjs/backend/knownWords';
import { backend } from '@wailsjs/models';
import DeleteLearned from '@/components/DeleteLearned.vue';
import ImportCsv from '@/components/ImportCsv.vue';
import {
  ImportAnkiKeywords,
  ImportAnkiReviewData,
} from '@wailsjs/backend/ankiInterface';
import { useLoader } from '@/lib/loading';
import { getUserSettings } from '@/lib/userSettings';
const UserSettings = getUserSettings();

const getRowId = (params:GetRowIdParams) => params.data.Word;

const loader = useLoader();
async function importAnki() {
  return loader.withLoader(ImportAnkiKeywords, false);
}

async function importAnkiReviews() {
  return loader.withLoader(ImportAnkiReviewData, false);
}

const columnDefs:ColDef[] = [
  {
    headerName: 'Word',
    field: 'Word',
    sortable: true,
    suppressSizeToFit: true,
    cellClass: 'text-xl',
  },
  {
    headerName: 'Interval',
    field: 'Interval',
    sort: 'desc',
    sortable: true,
    comparator: (_a, _b, rowA, rowB) => {
      return rowA.data.Interval - rowB.data.Interval;
    },
    cellClass: 'text-xl',
    valueGetter: (params: ValueGetterParams<backend.WordGridRow>) => {
      const interval = params.data?.Interval;
      if (interval === 10000) {
        return 'Marked Known';
      }
      return interval;
    },
  },
  {
    headerName: 'Learned On',
    field: 'LearnedOn',
    sortable: true,
    cellClass: 'text-xl',
  },
  {
    headerName: 'In Dicts',
    field: 'InDict',
    suppressSizeToFit: true,
    sortable: true,
    cellClass: 'text-xl',
  },
  {
    headerName: '',
    field: 'markButton',
    width: 120,
    cellRenderer: DeleteLearned,
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

const rowData = ref<backend.WordGridRow[]>([]);

onBeforeMount(() => {
  loadData();
});

onUnmounted(() => {
  window.removeEventListener('resize', resizeCallback);
});

async function loadData() {
  rowData.value = await GetWordsGrid();
}
</script>
