<template>
  <div v-if="store.wordList.length != 0" class="modal">
    <p>These are the words we are working on!
    {{store.wordList.join(" ")}}</p>

    <div v-if="step == 0">
      <p>Select the sentence</p>
      <div v-if="sentences.length > 0">
        <p v-for="(sentence, i) in sentences" :key="i">
        {{sentence}}
        </p>
      </div>
    </div>
    <div v-else-if="step == 1">
      <p>Select the Picture</p>
    </div>
    <div v-else-if="step == 2">
      <p>Select the Definition</p>
    </div>

      <n-button type=primary v-if="step>0" @click="step--">Previous</n-button>
      <n-button type=warning v-if="step<2" @click="step++">Next</n-button>
      <n-button type=info v-if="step==2">Submit</n-button>
    <n-button @click="store.wordList = []">Exit</n-button>
  </div>
</template>

<script>
import { defineComponent, ref } from 'vue';
import { useCardQueue } from '@/stores/CardQueue';

export default defineComponent({
  name: 'CardCreator',
  components: [
  ],
  setup() {
    const store = useCardQueue();
    const sentences = ref([]);
    store.$subscribe(async (mutation, state) => {
      // Later we can prefetch new words sentences possibly
      if (mutation.events.type === 'add' && mutation.events.key === '0') {
        const word = state.wordList[0];
        sentences.value = await window.ipc.getSentencesForWord(word);
      }
    });
    return {
      store, sentences,
    };
  },
  data() {
    return { step: 0 };
  },
});

</script>

<style scoped>
.modal {
  position: fixed;
  z-index: 999;
  top: 20%;
  left: 50%;
  width: 560px;
  margin-left: -150px;
  background: pink;
}
</style>
