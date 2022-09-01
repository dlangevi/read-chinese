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
        <p v-if="step == 1">Select the sentence</p>
        <p v-else-if="step==2">Select the Picture</p>
        <p v-else-if="step==3">Select the Definition</p>
      </template>
      <n-layout has-sider sider-placement="right">
        <n-layout-content content-style="padding: 24px;">
          <div v-if="step == 1">
            <n-checkbox-group v-if="sentences.length > 0" v-model:value="sentence">
              <n-space vertical item-style="display: flex;">
              <n-checkbox class="text-3xl" v-for="(sentence, i) in sentences" :key="i"
                :value="sentence" :label="sentence"/>
              </n-space>
            </n-checkbox-group>
          </div>
        </n-layout-content>
        <n-layout-sider v-if="card !== undefined"
          collapse-mode="transform"
          :collapsed-width="50"
          :native-scrollbar="true"
          :show-collapsed-content="false"
          :width="500"
          show-trigger="arrow-circle"
          content-style="padding: 24px;"
          bordered
        >
          <anki-card-preview :ankiCard="card"/>
        </n-layout-sider>
      </n-layout>

      <template #action>
        <n-space justify="end">
          <n-button type=primary v-if="step>1" @click="onNegativeClick">Previous</n-button>
          <n-button type=warning v-if="step<steps" @click="onPositiveClick">Next</n-button>
          <n-button type=info v-if="step==steps" @click="submit()">Submit</n-button>
        </n-space>
      </template>
    </n-modal>
</template>

<script setup>
import { ref, watch, reactive } from 'vue';
import { useCardQueue } from '@/stores/CardQueue';
import AnkiCardPreview from '@/components/AnkiCardPreview.vue';
import {
  useMessage, NCheckboxGroup, NCheckbox, NSpace, NButton, NModal,
  NLayoutSider, NLayout, NLayoutContent,
} from 'naive-ui';

const store = useCardQueue();
const sentences = ref([]);
const sentence = ref(null);
const step = ref(1);
const steps = 1;
const showModal = ref(false);
const currentWord = ref(null);
const card = ref(undefined);

const message = useMessage();
store.$subscribe(async (mutation, state) => {
  // Later we can prefetch new words sentences possibly
  if (mutation.events.type === 'add' && mutation.events.key === '0') {
    const word = state.wordList[0];
    sentences.value = await window.ipc.getSentencesForWord(word);
    // Todo card may not exist. In which case start a new one
    const ankiCard = await window.ipc.getAnkiCard(word);
    card.value = reactive(ankiCard);
  }
  step.value = 1;
  showModal.value = state.wordList.length !== 0;
  [currentWord.value] = state.wordList;
});

watch(sentence, (newSentence) => {
  console.log(newSentence);
  if (newSentence.length > 0) {
    const [s] = newSentence;
    card.value.fields.ExampleSentence.value = s;
    console.log('changed card');
  }
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
  message.info(JSON.stringify(card.value));
}

</script>
