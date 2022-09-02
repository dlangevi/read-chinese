<template>
  <n-button type="primary"
    @click="addToQueue">
    Replace Sentence
  </n-button>
</template>

<script setup>
import { useCardQueue } from '@/stores/CardQueue';
import { NButton } from 'naive-ui';

const props = defineProps({
  params: {
    type: Object,
  },
});

const store = useCardQueue();
async function addToQueue() {
  store.addWord(props.params.data.word);
}

</script>
