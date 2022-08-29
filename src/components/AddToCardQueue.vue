<template>
  <img @click="addToQueue(this.params.data)" class="
    font-black
    bg-red-600
    hover:bg-green-600
    absolute
    top-1/2
    rounded-xl
    -translate-y-1/2
    h-5
    w-5" src="../assets/circle-check.svg" alt="checkmark"/>
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
