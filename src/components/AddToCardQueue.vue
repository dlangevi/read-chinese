<template>
  <n-button type="primary"
    @click="addToQueue">
    {{params.text}}
  </n-button>
</template>

<script setup>
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import { NButton } from 'naive-ui';

const props = defineProps({
  params: {
    type: Object,
  },
});

const store = useCardQueue();
async function addToQueue() {
  const rowData = props.params.data;
  let action = ActionsEnum.MODIFY;
  if (props.params.create) {
    action = ActionsEnum.CREATE;
  }
  store.addWord(rowData.word, action, () => {
    props.params.api.applyTransaction({
      remove: [rowData],
    });
  });
}

</script>
