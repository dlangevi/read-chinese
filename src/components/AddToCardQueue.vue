<template>
  <n-button type="primary"
    @click="addToQueue(this.params.data)">
    Create new FlashCard
  </n-button>
</template>

<script>
import { useCardQueue } from '@/stores/CardQueue';

export default {
  data() {
    return {
    };
  },
  setup() {
    const store = useCardQueue();
    return {
      store,
    };
  },
  methods: {
    addToQueue(rowData) {
      this.store.addWord(rowData.word);
    },

  },
};
</script>
