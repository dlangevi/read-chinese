<template>
  <n-button type="primary"
    @click="markKnown">
    Mark Known
  </n-button>
</template>

<script setup>
import { NButton } from 'naive-ui';

const props = defineProps({
  params: {
    type: Object,
  },
});

function markKnown() {
  const rowData = props.params.data;
  // Keep with the convention of 10000 == user has claimed
  // they super known this one
  window.ipc.addWord(rowData.word, 10000);
  props.params.api.applyTransaction({
    remove: [rowData],
  });
}

</script>
