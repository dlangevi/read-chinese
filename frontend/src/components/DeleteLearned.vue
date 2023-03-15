<template>
  <div class="flex h-full items-center">
    <button
      class="btn-error btn-xs btn my-auto ml-auto"
      @click="deleteKnown"
    >
      Mark Unknown
    </button>
  </div>
</template>

<script lang="ts" setup>
import type { ICellRendererParams } from 'ag-grid-community';
import { DeleteWord } from '@wailsjs/backend/knownWords';
import {
  UpdateSentenceTable,
} from '@wailsjs/backend/Generator';
import { backend } from '@wailsjs/models';

const props = defineProps<{ params:ICellRendererParams<backend.WordGridRow>}>();

async function deleteKnown() {
  const rowData = props.params.data;
  if (!rowData) {
    return;
  }
  await DeleteWord(rowData.Word);
  props.params.api.applyTransaction({
    remove: [rowData],
  });

  UpdateSentenceTable(rowData.Word);
}
</script>
