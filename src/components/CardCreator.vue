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
          Creating card for {{card.fields.word}}
        </p>
      </template>
      <n-layout has-sider sider-placement="right" style="height: 500px">
        <n-layout-content content-style="padding: 24px;"
          :native-scrollbar="false">
          <edit-sentence v-if="step==StepsEnum.SENTENCE"
            :word="card.fields.word" :sentence="card.fields.sentence"
            @updateSentence="updateSentence"/>
          <edit-english-definition v-if="step==StepsEnum.ENGLISH"
            :word="card.fields.word" :definition="card.fields.englishDefn"
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
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
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
const showModal = ref(false);
const card = ref(undefined);
let originalValues;
const step = ref(undefined);
let word;
let action;
let callback;

const changeStep = (estep) => { step.value = estep; };
const message = useMessage();
store.$subscribe(async (mutation, state) => {
  // Later we can prefetch new words sentences possibly
  // if (mutation.events.type === 'add' && mutation.events.key === '0') {
  if (state.wordList.length > 0) {
    [{
      word,
      action,
      callback,
    }] = state.wordList;

    let ankiCard;
    if (action === ActionsEnum.CREATE) {
      ankiCard = await window.ipc.createAnkiNoteSkeleton(word);
    } else {
      ankiCard = await window.ipc.getAnkiNote(word);
    }
    card.value = reactive(ankiCard);
    originalValues = {
      ...ankiCard.fields,
    };
    step.value = 1;
    console.log(card.value);
    console.log(originalValues);
  }
  showModal.value = state.wordList.length !== 0;
});

const updateSentence = (newSentence) => {
  if (newSentence.length > 0) {
    card.value.fields.sentence = newSentence;
  }
};

const updateDefinition = (newDefinition) => {
  if (newDefinition.definition.length > 0) {
    card.value.fields.englishDefn = newDefinition.definition;
    card.value.fields.pinyin = newDefinition.pronunciation;
  }
};

function onClose() {
  store.clearWords();
  return false;
}

async function submit() {
  // Todo track changes to the card and submit those for update
  message.info('Card submited');
  // TODO figure out the logic for determining changes better

  if (action === ActionsEnum.CREATE) {
    console.log('creating card', card.value.fields);
    const res = await window.ipc.createAnkiCard({ ...card.value.fields });
    message.info(JSON.stringify(res));
  } else {
    const newData = {};
    Object.entries(card.value.fields).forEach(([field, value]) => {
      if (value !== originalValues[field]) {
        newData[field] = value;
      }
    });

    // Clear audio as sentence has changed
    if (newData.ExampleSentence) {
      newData.SentenceAudio = '';
    }
    const res = await window.ipc.updateAnkiCard(card.value.noteId, newData);
    message.info(JSON.stringify(res));
  }
  callback();
  store.clearFront();
}

</script>
