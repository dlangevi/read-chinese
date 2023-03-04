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

const props = defineProps<{
  params: ICellRendererParams,
}>();

const store = useCardQueue();
async function addToQueue() {
  const rowData = props.params.data;
  // For now while there are mixed columns
  const word = rowData.word || rowData.Word;
  store.addWord({ word },
    () => {
      props.params.api.applyTransaction({
        remove: [rowData],
      });
    },
    props.params.context?.bookId);
}

</script>
