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
    @click="addToQueue">
    {{params.text}}
  </button>
</template>

<script setup>
/*
import { NButton } from 'naive-ui';
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

const props = defineProps({
  params: {
    type: Object,
    required: true,
  },
});

const store = useCardQueue();
async function addToQueue() {
  const rowData = props.params.data;
  let action = ActionsEnum.MODIFY;
  if (props.params.create) {
    action = ActionsEnum.CREATE;
  }
  const options = {
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
