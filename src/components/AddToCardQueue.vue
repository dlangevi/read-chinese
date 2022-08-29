<template>
  <button type="button" class="inline-block px-4 py-1.5 bg-blue-600 text-white font-medium text-xs
    leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700
    focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition
  duration-150 ease-in-out"
    @click="addToQueue(this.params.data)">
    Create new FlashCard
  </button>
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
