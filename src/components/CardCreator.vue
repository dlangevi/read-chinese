<template>
    <n-modal
      class="w-4/5 max-h-[80vh]"
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
      </template>
      <n-layout has-sider sider-placement="right" style="height: 500px">
        <n-layout-content content-style="padding: 24px;"
          :native-scrollbar="false">
          <edit-sentence v-if="step==StepsEnum.SENTENCE"
            :word="word" :sentence="sentence"
            @updateSentence="updateSentence"/>
          <edit-english-definition v-if="step==StepsEnum.ENGLISH"
            :word="word" :definition="definition"
            @updateDefinition="updateDefinition"/>
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
          <anki-card-preview :ankiCard="card" @changeStep="changeStep"/>
        </n-layout-sider>
      </n-layout>

      <template #action>
        <n-space justify="end">
          <n-button type=info
            @click="submit()">Submit</n-button>
        </n-space>
      </template>
    </n-modal>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useCardQueue } from '@/stores/CardQueue';
import AnkiCardPreview from '@/components/AnkiCardPreview.vue';
import {
  useMessage, NSpace, NButton, NModal,
  NLayoutSider, NLayout, NLayoutContent,
} from 'naive-ui';
import EditSentence from '@/components/CardCreatorSteps/EditSentence.vue';
import EditEnglishDefinition from
  '@/components/CardCreatorSteps/EditEnglishDefinition.vue';
import StepsEnum from '@/components/CardCreatorSteps/StepsEnum';

const store = useCardQueue();
const sentence = ref(null);
const definition = ref(undefined);
const showModal = ref(false);
const currentWord = ref(null);
const card = ref(undefined);
const word = ref('');
const step = ref(undefined);

const changeStep = (estep) => { step.value = estep; };
const message = useMessage();
store.$subscribe(async (mutation, state) => {
  // Later we can prefetch new words sentences possibly
  // if (mutation.events.type === 'add' && mutation.events.key === '0') {
  if (state.wordList.length > 0) {
    [word.value] = state.wordList;
    // sentences.value = await window.ipc.getSentencesForWord(word);
    // console.log(sentences.value);
    // Todo card may not exist. In which case start a new one
    const ankiCard = await window.ipc.getAnkiNote(word.value);
    card.value = reactive(ankiCard);
    sentence.value = card.value.fields.ExampleSentence.value;
  }
  step.value = 1;
  showModal.value = state.wordList.length !== 0;
  [currentWord.value] = state.wordList;
});

const updateSentence = (newSentence) => {
  sentence.value = newSentence;
  if (newSentence.length > 0) {
    card.value.fields.ExampleSentence.value = newSentence;
  }
};

const updateDefinition = (newDefinition) => {
  definition.value = newDefinition;
  if (definition.value.length > 0) {
    card.value.fields.EnglishDefinition.value = newDefinition;
  }
};

function onClose() {
  store.clearWords();
  return false;
}

async function submit() {
  // Todo track changes to the card and submit those for update
  onClose();
  message.info('Card submited');
  // TODO figure out the logic for determining changes better
  const newData = {};
  if (sentence.value) {
    newData.ExampleSentence = sentence.value;
    // Since this is a new sentence, make sure to strip the previous audio
    newData.SentenceAudio = '';
  }
  if (definition.value) {
    newData.EnglishDefinition = definition.value;
  }
  const res = await window.ipc.updateAnkiCard(card.value.noteId, newData);
  message.info(JSON.stringify(res));
}

</script>
