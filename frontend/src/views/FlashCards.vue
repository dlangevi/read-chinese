<template>
  <with-sidebar>
    <template #sidebar>
      <div>
        <h3 class="text-lg font-semibold">Useful searches: </h3>
        <div class="m-4 flex flex-col gap-4">
          <button
            v-for="search,key in searches"
            :key="key"
            class="btn-primary btn-sm btn"
            @click="setSearch(search)"
          >
            Load {{ key }}
          </button>
        </div>
      </div>
      <div class="border-2 p-2 text-center">
        Found {{ rowData.length }} cards
      </div>
    </template>
    <div class="flex h-full flex-col gap-4 p-4">
      <div class="mx-auto flex w-5/6 flex-col">
        <div class="input-group">
          <input
            v-model="currentSearch"
            type="text"
            prefix="foo"
            placeholder="Input a search as if this were the anki browse window"
            class="input-primary input w-full"
          >
          <button class="btn" @click="doSearch()">Search</button>
        </div>
        <div
          class="ml-auto before:ml-0.5
        before:text-red-500 before:content-['*']"
        >
          Note all searches will be prefixed with
          <span class="font-mono text-secondary">{{ searchPrefix }}</span>
        </div>
      </div>
      <ag-grid-vue
        class="ag-theme-alpine mx-auto
        w-5/6 grow text-xl"
        :get-row-id="getRowId"
        :column-defs="columnDefs"
        :row-data="rowData"
        @grid-ready="onGridReady"
      />
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { AgGridVue } from 'ag-grid-vue3';
import { ref, onUnmounted } from 'vue';
import {
  getUserSettings,
} from '@/lib/userSettings';
import ProblemResolver from '@/components/ProblemResolver.vue';
import type {
  GetRowIdParams,
  GridReadyEvent, ColDef, ICellRendererParams,
} from 'ag-grid-community';
import { LoadProblemCards } from '@wailsjs/backend/ankiInterface';
import { backend } from '@wailsjs/models';
import WithSidebar from '@/layouts/WithSidebar.vue';

const UserSettings = getUserSettings();
const AnkiConfig = UserSettings.AnkiConfig;
const getRowId = (params:GetRowIdParams) => params.data.Word;

const activeModel = AnkiConfig.ActiveModel;
const currentMapping = AnkiConfig.ModelMappings[AnkiConfig.ActiveModel];
const activeDeck = AnkiConfig.ActiveDeck;
const currentSearch = ref('');

const searchPrefix = `deck:${activeDeck} note:${activeModel}`;
const searches = {
  Flagged: '-flag:0',
  Suspended: 'is:suspended',
  'Missing Sentence': `"${currentMapping.exampleSentence}:"`,
  'Missing Sentence Audio': `-"${
    currentMapping.exampleSentence
  }:" "${
    currentMapping.sentenceAudio
  }:"`,
  'Missing Image': `"${currentMapping.images}:"`,
  'Missing Hanzi Audio': `"${currentMapping.hanziAudio}:"`,
  'Missing Pinyin': `"${currentMapping.pinyin}:"`,
  'Long Sentence': `"${currentMapping.exampleSentence}:re:^.{40,}$"`,
};

const columnDefs:ColDef[] = [
  {
    headerName: 'word',
    field: 'Word',
    suppressSizeToFit: true,
    cellClass: 'text-xl',
  },
  {
    headerName: 'Detected Problems',
    field: 'Problems',
    sort: 'desc',
    cellClass: 'text-xl',
    autoHeight: true,
    cellRenderer: (params:ICellRendererParams) => {
      const issues = Object.entries(params.value)
        .filter(([_, value]) => { return value; })
        .map(([key, _]) => {
          return key;
        });
      if (params.data.Notes) {
        issues.push(`UserNote: ${params.data.Notes}`);
      }
      return `<ul>
      ${issues.map(issue => { return '<li>' + issue + '</li>'; }).join('')}
      </ul>`;
    },
  },
  {
    headerName: '',
    field: 'problemResolver',
    width: 200,
    cellClass: 'flex flex-row-reverse items-center',
    cellRenderer: ProblemResolver,
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

onUnmounted(() => {
  window.removeEventListener('resize', resizeCallback);
});

const rowData = ref<backend.ProblemCard[]>([]);
async function setSearch(search : string) {
  currentSearch.value = search;
  doSearch();
}

async function doSearch() {
  rowData.value = await LoadProblemCards(currentSearch.value);
}
</script>
