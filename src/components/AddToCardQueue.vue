<template>
  <n-button type="primary"
    @click="addToQueue(this.params.data)">
    Replace Sentence
  </n-button>
</template>

<script setup>
import { useCardQueue } from '@/stores/CardQueue';
import { NButton } from 'naive-ui';

const store = useCardQueue();

async function addToQueue(rowData) {
  store.addWord(rowData.word);
}

</script>
