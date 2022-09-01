<template>
    <n-modal
      class="w-4/5"
      v-model:show="showModal"
      :mask-closable="false"
      :closable="true"
      preset="card"
      :on-close="onClose"
    >
      <template #header>
        <p class="text-xl">
          Creating card for {{currentWord}}
        </p>
        <p v-if="step == 0">Select the sentence</p>
        <p v-else-if="step==1">Select the Picture</p>
        <p v-else-if="step==2">Select the Definition</p>
      </template>
      <div v-if="step == 0">
        <n-checkbox-group v-if="sentences.length > 0" v-model:value="sentence">
          <n-space vertical item-style="display: flex;">
          <n-checkbox v-for="(sentence, i) in sentences" :key="i"
            :value="sentence" :label="sentence"/>
          </n-space>
        </n-checkbox-group>
      </div>
      <div v-else-if="step == 1">
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
        <n-checkbox-group v-model:value="food">
          <n-space vertical item-style="display: flex;">
            <n-checkbox value="Potatos" label="Potatos" />
            <n-checkbox value="Fish" label="Fish" />
          </n-space>
        </n-checkbox-group>
      </div>
      <template #action>
        <n-space justify="end">
          <n-button type=primary v-if="step>0" @click="onNegativeClick">Previous</n-button>
          <n-button type=warning v-if="step<2" @click="onPositiveClick">Next</n-button>
          <n-button type=info v-if="step==2" @click="submit()">Submit</n-button>
        </n-space>
      </template>
    </n-modal>
</template>

<script setup>
import { ref } from 'vue';
import { useCardQueue } from '@/stores/CardQueue';
import {
  useMessage, NCheckboxGroup, NCheckbox, NSpace, NButton, NModal,
} from 'naive-ui';

const store = useCardQueue();
const sentences = ref([]);
const sentence = ref(null);
const cities = ref(null);
const food = ref(null);
const step = ref(0);
const showModal = ref(false);
const currentWord = ref(null);

const message = useMessage();
store.$subscribe(async (mutation, state) => {
  // Later we can prefetch new words sentences possibly
  if (mutation.events.type === 'add' && mutation.events.key === '0') {
    const word = state.wordList[0];
    sentences.value = await window.ipc.getSentencesForWord(word);
  }
  step.value = 0;
  showModal.value = state.wordList.length !== 0;
  [currentWord.value] = state.wordList;
});

function onNegativeClick() {
  step.value -= 1;
  return false;
}
function onPositiveClick() {
  step.value += 1;
  return false;
}
function onClose() {
  store.clearWords();
  return false;
}

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
