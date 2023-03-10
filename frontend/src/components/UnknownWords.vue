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
import type { GetRowIdParams, GridReadyEvent, ColDef } from 'ag-grid-community';
import { GetDefinitions } from '@wailsjs/backend/Dictionaries';
import {
  GetBookFrequencies,
  GetFavoriteFrequencies,
} from '@wailsjs/backend/bookLibrary';
import { GetOccurances } from '@wailsjs/backend/KnownWords';
import type { WordDefinitions, UnknownWordRow } from '@/lib/types';

import { useCardQueue } from '@/stores/CardQueue';

const UserSettings = getUserSettings();

const props = defineProps<{
  words: string[],
  bookFilter?: number,
  frequencySource?: string,
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
  ...(UserSettings.Dictionaries.ShowDefinitions
    ? [
      {
        headerName: 'definition',
        field: 'definition',
        minWidth: 400,
      }]
    : []),
  {
    headerName: '',
    field: 'markButton',
    width: 120,
    cellRenderer: MarkLearned,
  },
  {
    headerName: '',
    field: 'Make FlashCard',
    width: 120,
    cellRenderer: AddToCardQueue,
  },
];

const rowData = ref<UnknownWordRow[]>([]);
watch(() => props.words, async () => {
  console.log(props.words);
  updateWords();
});

async function updateWords() {
  console.log(props.words);
  const definitions : WordDefinitions = await GetDefinitions(props.words);
  let occurances : {
    [key:string] :number
  } = {};
  if (props.bookFilter) {
    occurances = await GetBookFrequencies(props.bookFilter);
  } else if (props.frequencySource === 'favorites') {
    occurances = await GetFavoriteFrequencies();
  } else {
    occurances = await GetOccurances(props.words);
  }
  rowData.value = props.words.map((word) => {
    const row : UnknownWordRow = { word };
    const definition = definitions[word];
    if (definition) {
      row.definition = definition.definition;
      row.pinyin = definition.pronunciation;
    }
    row.occurance = occurances[word];
    return row;
  });
}

const store = useCardQueue();
defineExpose({
  enqueueTopRows: async function (n : number) {
    const sorted = rowData.value.slice();
    sorted.sort((a, b) => {
      if (a.occurance === undefined || b.occurance === undefined) {
        return 0;
      }
      if (a.occurance > b.occurance) {
        return -1;
      }
      return 1;
    });
    const topWords = sorted.slice(0, n);
    topWords.forEach((word) => {
      store.addWord({ word: word.word });
    });
  },
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
  updateWords();
});

</script>
