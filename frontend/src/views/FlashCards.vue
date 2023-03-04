<template>
  <with-sidebar>
    <template #sidebar>
      <div>
        <p>Filtered problem cards: </p>
        <div>
          <settings-checkbox
            v-for="([initial, setting]) in
              getDisplayable(UserSettings.CardManagement)"
            :key="setting.name"
            :setting="setting"
            :initial-value="initial"
          />
        </div>
      </div>
      <div class="border-2 p-2 text-center">
        {{ rowData.length }} total 'problem' cards
      </div>
      <div class="border-2 p-2 text-center">
        {{ visibleRows }} cards after filter
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
        @filter-changed="filterChanged"
        @grid-ready="onGridReady"
      />
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { AgGridVue } from 'ag-grid-vue3';
import { onBeforeMount, ref, onUnmounted, watch } from 'vue';
import {
  getDisplayable,
  getUserSettings,
} from '@/lib/userSettings';
import ProblemResolver from '@/components/ProblemResolver.vue';
import type {
  GetRowIdParams, GridApi,
  GridReadyEvent, ColDef, ICellRendererParams, RowNode,
} from 'ag-grid-community';
import { LoadProblemCards } from '@wailsjs/backend/ankiInterface';
import { backend } from '@wailsjs/models';
import WithSidebar from '@/layouts/WithSidebar.vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';

const UserSettings = getUserSettings();
const CardManagement = UserSettings.CardManagement;
const getRowId = (params:GetRowIdParams) => params.data.Word;

// Will be set on grid ready
const gridApi = ref<GridApi>();
const visibleRows = ref(0);
function filterChanged() {
  if (gridApi.value) {
    visibleRows.value = gridApi.value.getModel().getRowCount();
  }
}
watch(() => ([
  CardManagement.ProblemFlagged,
  CardManagement.ProblemMissingImage,
  CardManagement.ProblemMissingSentence,
  CardManagement.ProblemMissingSentenceAudio,
  CardManagement.ProblemMissingWordAudio,
  CardManagement.ProblemMissingPinyin,
]), () => {
  gridApi.value?.onFilterChanged();
  if (gridApi.value) {
    visibleRows.value = gridApi.value.getModel().getRowCount();
  }
},

);

function isExternalFilterPresent() {
  return true;
}

function doesExternalFilterPass(node: RowNode) {
  const problems : backend.Problems = node.data.Problems;
  const ProblemFlagged = CardManagement.ProblemFlagged;
  const ProblemMissingImage = CardManagement.ProblemMissingImage;
  const ProblemMissingSentence = CardManagement.ProblemMissingSentence;
  const ProblemMissingSentenceAudio =
    CardManagement.ProblemMissingSentenceAudio;
  const ProblemMissingWordAudio = CardManagement.ProblemMissingWordAudio;
  const ProblemMissingPinyin = CardManagement.ProblemMissingPinyin;

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
    field: 'Fix the problem',
    width: 200,
    cellClass: 'flex flex-row-reverse items-center',
    cellRenderer: ProblemResolver,
  },
];

let resizeCallback: () => void;
function onGridReady(params:GridReadyEvent) {
  params.api.sizeColumnsToFit();
  gridApi.value = params.api;
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
  rowData.value = await LoadProblemCards();
  gridApi.value?.onFilterChanged();
});
</script>
