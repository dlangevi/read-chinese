<template>
  <button
    class="btn btn-xs btn-secondary"
    @click="markKnown"
  >
    Mark Known
  </button>
</template>

<script lang="ts" setup>
import type { ICellRendererParams } from 'ag-grid-community';
import { AddWord } from '@wailsjs/backend/KnownWords';

const props = defineProps<{ params:ICellRendererParams }>();

function markKnown() {
  const rowData = props.params.data;
  // Keep with the convention of 10000 == user has claimed to known a word
  AddWord(rowData.word, 10000);
  props.params.api.applyTransaction({
    remove: [rowData],
  });
}
</script>
