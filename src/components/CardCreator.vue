<template>
  <div v-if="store.wordList.length != 0" class="modal">
    <p>These are the words we are working on!</p>
    <p>{{wordList}}</p>
    <button @click="store.wordList = []">Clear</button>
  </div>
</template>

<script>
import { mapState } from 'pinia';
import { useCardQueue } from '@/stores/CardQueue';

export default {
  name: 'CardCreator',
  setup() {
    const store = useCardQueue();
    return {
      store,
    };
  },
  computed: {
    ...mapState(useCardQueue, ['wordList']),
  },
  data() {
    return { open: false };
  },
};

</script>

<style scoped>
.modal {
  position: fixed;
  z-index: 999;
  top: 20%;
  left: 50%;
  width: 300px;
  margin-left: -150px;
  background: pink;
}
</style>
