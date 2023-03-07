<template>
  <button
    class="btn-secondary btn-xs btn"
    @click="addToQueue"
  >
    Create FlashCard
  </button>
</template>

<script lang="ts" setup>
import { useCardQueue } from '@/stores/CardQueue';
import type { ICellRendererParams } from 'ag-grid-community';
import type { UnknownWordRow } from '@/lib/types';

const props = defineProps<{
  params: ICellRendererParams<UnknownWordRow>,
}>();

const store = useCardQueue();
async function addToQueue() {
  const rowData = props.params.data;
  if (rowData) {
    store.addWord({ word: rowData.word },
      () => {
        props.params.api.applyTransaction({
          remove: [rowData],
        });
      },
      props.params.context?.bookId);
  }
}

</script>
