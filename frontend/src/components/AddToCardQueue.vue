<template>
  <button
    type="button"
    class="btn text-sm"
    @click="addToQueue"
  >
    {{ params.text }}
  </button>
</template>

<script lang="ts" setup>
/*
import { NButton } from 'naive-ui';
    class="bg-purple-700
           text-white
           hover:bg-purple-500
           hover:text-white
           px-1
           rounded-lg
           text-sm"
  <n-button
    type="default"
    size="tiny"
    round
    color="#8a2be2"
    class="w-28"
    @click="addToQueue">
    {{params.text}}
  </n-button>

*/
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
