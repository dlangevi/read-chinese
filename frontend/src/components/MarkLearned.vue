<template>
  <button
    class="btn-secondary btn-xs btn"
    @click="markKnown"
  >
    Mark Known
  </button>
</template>

<script lang="ts" setup>
import type { ICellRendererParams } from 'ag-grid-community';
import { AddWord } from '@wailsjs/backend/KnownWords';
import {
  UpdateSentenceTable,
} from '@wailsjs/backend/Generator';

const props = defineProps<{ params:ICellRendererParams }>();

async function markKnown() {
  const rowData = props.params.data;
  // Keep with the convention of 10000 == user has claimed to known a word
  await AddWord(rowData.word, 10000);
  UpdateSentenceTable(rowData.word);
}
</script>
