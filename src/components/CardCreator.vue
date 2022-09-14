<template>
  <n-modal
    class="w-4/5 h-[80vh]"
    v-model:show="showModal"
    :mask-closable="false"
    :closable="true"
    preset="card"
    :on-close="onClose"
    content-style="display:flex; flex-direction: column;"
  >
    <template #header>
      <p class="text-xl">
        Creating card for {{card.fields.word}}
      </p>
    </template>
    <template #header-extra>
      <card-creation-settings/>
    </template>
    <n-layout has-sider sider-placement="left" class="flex-grow" >
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
      <n-layout-content content-style="padding: 24px;"
        :native-scrollbar="false">
        <edit-sentence v-if="step==StepsEnum.SENTENCE"
          :word="card.fields.word" :sentence="card.fields.sentence"
          @updateSentence="updateSentence"/>
        <edit-english-definition v-if="step==StepsEnum.ENGLISH"
          :word="card.fields.word" :definition="card.fields.englishDefn"
          @updateDefinition="updateDefinition"/>
      </n-layout-content>
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
import CardCreationSettings from '@/components/CardCreationSettings.vue';
import {
  useMessage, NSpace, NButton, NModal,
  NLayoutSider, NLayout, NLayoutContent,
} from 'naive-ui';
import EditSentence from '@/components/CardCreatorSteps/EditSentence.vue';
import EditEnglishDefinition from
  '@/components/CardCreatorSteps/EditEnglishDefinition.vue';
import StepsEnum from '@/components/CardCreatorSteps/StepsEnum';
import UserSettings from '@/userSettings';

const store = useCardQueue();
const showModal = ref(false);
const card = ref(undefined);
let originalValues;
const step = ref(undefined);
let steps = [];
let word;
let action;
let callback;

// Manually change the step from an edit button.
const changeStep = (estep) => {
  step.value = estep;
  // Remove any set steps progression since the user has taken control
  steps = [];
};

const nextStep = (currentStep) => {
  if (steps.length > 0 && currentStep === step.value) {
    const idx = steps.indexOf(currentStep);
    if (idx !== -1 && idx + 1 <= steps.length) {
      step.value = steps[idx + 1];
    }
  }
};

const updateSentence = (newSentence) => {
  if (newSentence.length > 0) {
    card.value.fields.sentence = newSentence;
    nextStep(StepsEnum.SENTENCE);
  }
};

const updateDefinition = (newDefinitions, updateStep = true) => {
  console.log(newDefinitions);
  if (newDefinitions.length > 0) {
    card.value.fields.englishDefn = newDefinitions.map(
      (def) => def.definition,
    ).join('<br>');
    card.value.fields.pinyin = newDefinitions.map(
      (def) => def.pronunciation,
    ).join(', ');
    if (updateStep) {
      nextStep(StepsEnum.ENGLISH);
    }
  }
};

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
      steps = [
        StepsEnum.SENTENCE,
        StepsEnum.ENGLISH,
      ];
    } else {
      // Right now for EDIT we only edit the sentence so start there
      ankiCard = await window.ipc.getAnkiNote(word);
      steps = [
        StepsEnum.SENTENCE,
      ];
    }
    card.value = reactive(ankiCard);
    originalValues = {
      ...ankiCard.fields,
    };

    // TODO this needs to be written in a more modular way
    const englishIdx = steps.indexOf(StepsEnum.ENGLISH);
    if (englishIdx !== -1) {
      const autoFill = await UserSettings.PopulateEnglish.read();
      if (autoFill) {
        const definitions = await window.ipc.getDefinitionsForWord(word);
        if (definitions.length === 1) {
          updateDefinition(definitions);
          steps.splice(englishIdx, 1);
        }
      }
    }

    [step.value] = steps;
    console.log(card.value);
    console.log(originalValues);
  }
  showModal.value = state.wordList.length !== 0;
});

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
