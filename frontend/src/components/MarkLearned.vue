<template>
  <button
    type="button"
    class="bg-purple-700
           text-white
           hover:bg-purple-500
           hover:text-white
           px-1
           rounded-lg
           text-sm"
    @click="markKnown"
  >
    Mark Known
  </button>
</template>

<script lang="ts" setup>
/*
import { NButton } from 'naive-ui';
  <n-button
    type="default"
    size="tiny"
    round
    color="#8a2be2"
    class="w-28"
    @click="markKnown">
    Mark Known
  </n-button>
*/
import type { ICellRendererParams } from 'ag-grid-community';
import { AddWord } from '@wailsjs/backend/KnownWords';

const props = defineProps<{ params:ICellRendererParams }>();

function markKnown() {
  const rowData = props.params.data;
  // Keep with the convention of 10000 == user has claimed to known a word
  console.log('adding with go');
  AddWord(rowData.word, 10000);
  props.params.api.applyTransaction({
    remove: [rowData],
  });
}
</script>
