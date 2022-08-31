<template>
  <n-button type="primary"
    @click="addToQueue(this.params.data)">
    Create new FlashCard
  </n-button>
</template>

<script setup>
import { useCardQueue } from '@/stores/CardQueue';

const store = useCardQueue();

function addToQueue(rowData) {
  store.addWord(rowData.word);
}

</script>
