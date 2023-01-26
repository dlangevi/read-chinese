<template>
  <with-sidebar>
    <template #sidebar>
      <div>
        <p>Total problem cards: {{ rowData.length }}</p>
        <p>Filtered problem cards: </p>
        <div>
          <component
            :is="content.type"
            v-for="(content) in CardManagement"
            :key="content"
            :setting="content"
            @update="updateFilter"
          />
        </div>
      </div>
    </template>
    <div class="container m-4 mx-auto flex h-full flex-col px-4">
      <ag-grid-vue
        class="ag-theme-alpine mx-auto h-full
        w-5/6 grow text-xl"
        :get-row-id="getRowId"
        :column-defs="columnDefs"
        :row-data="rowData"
        :is-external-filter-present="isExternalFilterPresent"
        :does-external-filter-pass="doesExternalFilterPass"
        @grid-ready="onGridReady"
      />
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { AgGridVue } from 'ag-grid-vue3';
import { onBeforeMount, ref, onUnmounted } from 'vue';
import { getUserSettings } from '@/lib/userSettings';
import { useMessage } from '@/lib/messages';
import AddToCardQueue from '@/components/AddToCardQueue.vue';
import type {
  GetRowIdParams, GridApi,
  GridReadyEvent, ColDef, ICellRendererParams, RowNode,
} from 'ag-grid-community';
import { LoadProblemCards } from '@wailsjs/backend/AnkiInterface';
import { backend } from '@wailsjs/models';
import WithSidebar from '@/layouts/WithSidebar.vue';

const UserSettings = getUserSettings();
const CardManagement = UserSettings.CardManagement;
const message = useMessage();
const getRowId = (params:GetRowIdParams) => params.data.Word;

let ProblemFlagged = CardManagement.ProblemFlagged.read();
let ProblemMissingImage = CardManagement.ProblemMissingImage.read();
let ProblemMissingSentence = CardManagement.ProblemMissingSentence.read();
let ProblemMissingSentenceAudio =
  CardManagement.ProblemMissingSentenceAudio.read();
let ProblemMissingWordAudio = CardManagement.ProblemMissingWordAudio.read();
let ProblemMissingPinyin = CardManagement.ProblemMissingPinyin.read();

// Will be set on grid ready
let gridApi : GridApi;
function updateFilter() {
  ProblemFlagged = CardManagement.ProblemFlagged.read();
  ProblemMissingImage = CardManagement.ProblemMissingImage.read();
  ProblemMissingSentence = CardManagement.ProblemMissingSentence.read();
  ProblemMissingSentenceAudio =
    CardManagement.ProblemMissingSentenceAudio.read();
  ProblemMissingWordAudio = CardManagement.ProblemMissingWordAudio.read();
  ProblemMissingPinyin = CardManagement.ProblemMissingPinyin.read();
  gridApi.onFilterChanged();
}

function isExternalFilterPresent() {
  return true;
}

function doesExternalFilterPass(node: RowNode) {
  const problems : backend.Problems = node.data.Problems;

  const passes = (ProblemFlagged && problems.Flagged) ||
     (ProblemMissingImage && problems.MissingImage) ||
     (ProblemMissingSentence && problems.MissingSentence) ||
     (ProblemMissingSentenceAudio && problems.MissingSentenceAudio) ||
     (ProblemMissingWordAudio && problems.MissingWordAudio) ||
     (ProblemMissingPinyin && problems.MissingPinyin);
  return passes;
}

const columnDefs:ColDef[] = [
  {
    headerName: 'word',
    field: 'Word',
    suppressSizeToFit: true,
    cellClass: 'text-xl',
  },
  {
    headerName: 'problem',
    field: 'Problems',
    sort: 'desc',
    cellClass: 'text-xl',
    autoHeight: true,
    cellRenderer: (params:ICellRendererParams) => {
      // put the value in bold
      console.log(params);
      const issues = Object.entries(params.value)
        .filter(([_, value]) => { return value; })
        .map(([key, _]) => {
          return key;
        });
      console.log(params.data);
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
    field: 'Make FlashCard',
    width: 50,
    cellRenderer: AddToCardQueue,
    cellRendererParams: {
      text: 'Replace Sentence',
    },
  },
];

let resizeCallback: () => void;
function onGridReady(params:GridReadyEvent) {
  params.api.sizeColumnsToFit();
  gridApi = params.api;
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
onBeforeMount(async () => {
  // TODO is there some way we can have the message be sent from
  // background, so we can just error at the source and not have
  // to do error handleing in all of the different vue files?
  try {
    // rowData.value = await LoadFlaggedCards();
    rowData.value = await LoadProblemCards();
  } catch (e) {
    message.error('Please open Anki');
  }
});
</script>
