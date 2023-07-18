<template>
  <ag-grid-vue
    class="ag-theme-alpine text-xl"
    :column-defs="columnDefs"
    row-model-type="infinite"
    @grid-ready="onGridReady"
  />
</template>

<script lang="ts" setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import {
  onBeforeMount, onUnmounted,
} from 'vue';
import { AgGridVue } from 'ag-grid-vue3';
import MarkLearned from '@/components/MarkLearned.vue';
import AddToCardQueue from '@/components/AddToCardQueue.vue';
import { getUserSettings } from '@/lib/userSettings';
import type {
  GridReadyEvent,
  ColDef, GridApi,
} from 'ag-grid-community';
import type { UnknownWordRow } from '@/lib/types';

import { useCardQueue } from '@/stores/CardQueue';

import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { InfiniteScroll } from '@/lib/infiniteScroll';

const UserSettings = getUserSettings();

const columnDefs:ColDef[] = [
  {
    headerName: 'word',
    field: 'word',
    // width: 80,
    // minWidth: 80,
    cellClass: 'text-xl',
  },
  {
    headerName: 'pinyin',
    field: 'pinyin',
    // width: 100,
    cellClass: [
      'border-2',
      'text-base-content',
      ...(UserSettings.Dictionaries.ShowPinyin
        ? [
        ]
        : [
          'text-opacity-0',
          'hover:text-opacity-100',
        ]),
    ],
  },
  {
    headerName: 'occurances',
    field: 'occurance',
    // width: 50,
    minWidth: 50,
  },
  {
    headerName: 'frequency',
    field: 'frequency',
    // width: 80,
    minWidth: 50,
    valueFormatter: function (params) {
      if (!params.data) {
        return 'none';
      }
      const rank = params.data.frequency;
      if (rank === undefined) {
        return 'none';
      } else if (rank < 1500) {
        return '★★★★★';
      } else if (rank < 5000) {
        return '★★★★';
      } else if (rank < 15000) {
        return '★★★';
      } else if (rank < 30000) {
        return '★★';
      } else if (rank < 60000) {
        return '★';
      } else {
        return 'none';
      }
    },
  },
  {
    headerName: 'definition',
    field: 'definition',
    minWidth: 400,
    cellClass: [
      'text-base-content',
      ...(UserSettings.Dictionaries.ShowDefinitions
        ? [
        ]
        : [
          'text-opacity-0',
          'hover:text-opacity-100',
        ]),
    ],
  },
  {
    headerName: '',
    field: 'markButton',
    // width: 120,
    cellRenderer: MarkLearned,
  },
  {
    headerName: '',
    field: 'Make FlashCard',
    // width: 120,
    cellRenderer: AddToCardQueue,
  },
];

let gridApi : GridApi<UnknownWordRow>;
let resizeCallback: () => void;
const dataSource = new InfiniteScroll();
const store = useCardQueue();

defineExpose({
  enqueueTopRows: async function (n : number) {
    const sortedRows : UnknownWordRow[] = [];
    gridApi.forEachNode((node) => {
      if (node.data) {
        sortedRows.push(node.data);
      }
    });
    const topWords = sortedRows.slice(0, n);
    topWords.forEach((word) => {
      store.addWord({ word: word.word });
    });
  },
  getRowCount: function () {
    return dataSource.rowCount;
  },
});

function onGridReady(params:GridReadyEvent) {
  gridApi = params.api;
  gridApi.setDatasource(dataSource);
  resizeCallback = () => {
    setTimeout(() => {
      params.api.sizeColumnsToFit();
    }, 10);
  };
  window.addEventListener('resize', resizeCallback);
  resizeCallback();
  EventsOn('AddedWord', (_ : string) => {
    params.api.purgeInfiniteCache();
    if (resizeCallback) {
      resizeCallback();
    }
  });

  EventsOn('ResetBoard', () => {
    console.log('Resetting Board');
    params.api.purgeInfiniteCache();
    if (resizeCallback) {
      resizeCallback();
    }
  });
}

onUnmounted(() => {
  window.removeEventListener('resize', resizeCallback);
  EventsOff('AddedWord');
  EventsOff('ResetBoard');
});

onBeforeMount(async () => {
  // updateWords();
});

</script>
