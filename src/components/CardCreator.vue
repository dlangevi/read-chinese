<template>
  <div v-if="store.wordList.length != 0" class="modal">
    <p>These are the words we are working on!
    {{store.wordList.join(" ")}}</p>

    <div v-if="step == 0">
      <p>Select the sentence</p>
      <n-checkbox-group v-if="sentences.length > 0" v-model:value="sentence">
        <n-space vertical item-style="display: flex;">
        <n-checkbox v-for="(sentence, i) in sentences" :key="i"
          :value="sentence" :label="sentence"/>
        </n-space>
      </n-checkbox-group>
    </div>
    <div v-else-if="step == 1">
      <p>Select the Picture</p>
      <n-checkbox-group v-model:value="cities">
        <n-space vertical item-style="display: flex;">
          <n-checkbox value="Beijing" label="Beijing" />
          <n-checkbox value="Shanghai" label="Shanghai" />
          <n-checkbox value="Guangzhou" label="Guangzhou" />
          <n-checkbox value="Shenzen" label="Shenzhen" />
        </n-space>
      </n-checkbox-group>
    </div>
    <div v-else-if="step == 2">
      <p>Select the Definition</p>
      <n-checkbox-group v-model:value="food">
        <n-space vertical item-style="display: flex;">
          <n-checkbox value="Potatos" label="Potatos" />
          <n-checkbox value="Fish" label="Fish" />
        </n-space>
      </n-checkbox-group>
    </div>

      <n-button type=primary v-if="step>0" @click="step--">Previous</n-button>
      <n-button type=warning v-if="step<2" @click="step++">Next</n-button>
      <n-button type=info v-if="step==2" @click="submit()">Submit</n-button>
    <n-button @click="store.wordList = []">Exit</n-button>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useCardQueue } from '@/stores/CardQueue';
import { useMessage } from 'naive-ui';

const store = useCardQueue();
const sentences = ref([]);
const sentence = ref(null);
const cities = ref(null);
const food = ref(null);
const step = ref(0);

const message = useMessage();
store.$subscribe(async (mutation, state) => {
  // Later we can prefetch new words sentences possibly
  if (mutation.events.type === 'add' && mutation.events.key === '0') {
    const word = state.wordList[0];
    sentences.value = await window.ipc.getSentencesForWord(word);
  }
});

function submit() {
  message.info(JSON.stringify([sentence, cities, food]));
  console.log(JSON.stringify([sentence, cities, food]));
}

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
  max-height: 800px;
}
</style>
