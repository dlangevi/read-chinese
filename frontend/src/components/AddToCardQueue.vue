<template>
  <button
    class="btn-secondary btn-xs btn"
    @click="addToQueue"
  >
    {{ params.text }}
  </button>
</template>

<script lang="ts" setup>
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import type { ICellRendererParams } from 'ag-grid-community';
import type { WordOptions } from '@/stores/CardQueue';

interface ButtonParams extends ICellRendererParams {
  text: string;
  create?: boolean;
}

const props = defineProps<{
  params: ButtonParams,
}>();

const store = useCardQueue();
async function addToQueue() {
  const rowData = props.params.data;
  let action: ActionsEnum = ActionsEnum.MODIFY;
  if (props.params.create) {
    action = ActionsEnum.CREATE;
  }
  const options:WordOptions = {
    callback: () => {
      props.params.api.applyTransaction({
        remove: [rowData],
      });
    },
  };
  const { context } = props.params;
  if (context !== undefined) {
    options.preferBook = props.params.context.bookId;
  }
  store.addWord(rowData.word, action, options);
}

</script>
