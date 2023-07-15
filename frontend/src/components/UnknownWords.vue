<template>
  <ag-grid-vue
    class="ag-theme-alpine text-xl"
    :get-row-id="getRowId"
    :column-defs="columnDefs"
    :row-data="rowData"
    :context="gridContext"
    @grid-ready="onGridReady"
  />
</template>

<script lang="ts" setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import {
  watch, onBeforeMount, ref, onUnmounted,
} from 'vue';
import { AgGridVue } from 'ag-grid-vue3';
import MarkLearned from '@/components/MarkLearned.vue';
import AddToCardQueue from '@/components/AddToCardQueue.vue';
import { getUserSettings } from '@/lib/userSettings';
import type {
  GetRowIdParams, GridReadyEvent,
  ColDef, GridApi,
} from 'ag-grid-community';
import {
  GetWordData,
} from '@wailsjs/backend/wordLists';
import type { UnknownWordRow } from '@/lib/types';

import { useCardQueue } from '@/stores/CardQueue';

import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

const UserSettings = getUserSettings();

const props = withDefaults(defineProps<{
  words: string[],
  bookFilter?: number,
  sortByFrequency?: boolean,
  occuranceSource?: string,
  frequencySource?: string,
}>(), {
  bookFilter: undefined,
  sortByFrequency: false,
  occuranceSource: 'all',
  frequencySource: 'Jieba',
});

const getRowId = (params:GetRowIdParams) => params.data.word;

const gridContext = {
  bookId: props.bookFilter,
};

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
    sort: 'desc',
    sortIndex: props.sortByFrequency ? 1 : 0,
    // width: 50,
    minWidth: 50,
  },
  {
    headerName: 'frequency',
    field: 'frequency',
    sort: 'asc',
    sortIndex: props.sortByFrequency ? 0 : 1,
    // width: 80,
    minWidth: 50,
    valueFormatter: function (params) {
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

let resizeCallback: () => void;
let updateSorts: () => void;
const rowData = ref<UnknownWordRow[]>([]);
watch(() => [
  props.words,
  props.frequencySource,
  props.occuranceSource,
], async () => {
  updateWords();
});

watch(() => props.sortByFrequency, async () => {
  updateSorts();
});

async function updateWords() {
  let occuranceSource = props.occuranceSource;
  if (props.bookFilter) {
    occuranceSource = props.bookFilter.toString();
  }
  rowData.value = await GetWordData(
    props.words,
    occuranceSource,
    props.frequencySource);
  if (resizeCallback) {
    resizeCallback();
  }
}

let gridApi : GridApi<UnknownWordRow>;
const store = useCardQueue();
defineExpose({
  enqueueTopRows: async function (n : number) {
    const sortedRows : UnknownWordRow[] = [];
    gridApi.forEachNodeAfterFilterAndSort((node) => {
      if (node.data) {
        sortedRows.push(node.data);
      }
    });
    const topWords = sortedRows.slice(0, n);
    topWords.forEach((word) => {
      store.addWord({ word: word.word });
    });
  },
});

function onGridReady(params:GridReadyEvent) {
  gridApi = params.api;
  resizeCallback = () => {
    setTimeout(() => {
      params.api.sizeColumnsToFit();
    }, 10);
  };
  updateSorts = () => {
    params.columnApi.applyColumnState({
      state: [{
        colId: 'occurance',
        sortIndex: props.sortByFrequency ? 1 : 0,
      },
      {
        colId: 'frequency',
        sortIndex: props.sortByFrequency ? 0 : 1,
      },
      ],
    });
  };
  window.addEventListener('resize', resizeCallback);
  resizeCallback();
  EventsOn('AddedWord', (word : string) => {
    // TODO this might end up being slow, or could have
    // problems if markknown is spam clicked
    rowData.value = rowData.value.filter(
      (row) => row.word !== word);
  });
}

onUnmounted(() => {
  window.removeEventListener('resize', resizeCallback);
  EventsOff('AddedWord');
});

onBeforeMount(async () => {
  updateWords();
});
</script>
