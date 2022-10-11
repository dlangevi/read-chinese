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
      <card-creation-settings />
    </template>
    <n-layout
      has-sider
      sider-placement="left"
      class="flex-grow"
      style="max-height: 60vh">
      <n-layout-sider
        v-if="card !== undefined"
        collapse-mode="transform"
        :collapsed-width="50"
        :show-collapsed-content="false"
        :width="500"
        show-trigger="arrow-circle"
        content-style="padding: 24px;"
        bordered
      >
        <Suspense>
          <anki-card-preview :ankiCard="card" @changeStep="changeStep" />
        </Suspense>
      </n-layout-sider>
      <n-layout-content
        content-style="padding: 24px;"
        :native-scrollbar="true">
        <edit-sentence
          v-if="step == StepsEnum.SENTENCE"
          :preferBook="preferBookRef"
          :word="card.fields.word"
          :sentence="card.fields.sentence"
          @updateSentence="updateSentence" />
        <edit-definition
          v-if="step == StepsEnum.ENGLISH"
          :word="card.fields.word"
          :definition="card.fields.englishDefn"
          type="english"
          @updateDefinition="updateEnglishDefinition" />
        <edit-definition
          v-if="step == StepsEnum.CHINESE"
          :word="card.fields.word"
          :definition="card.fields.englishDefn"
          type="chinese"
          @updateDefinition="updateChineseDefinition" />
        <edit-images
          v-if="step == StepsEnum.IMAGE"
          :word="card.fields.word"
          @updateImages="updateImages" />
      </n-layout-content>
    </n-layout>

    <template #action>
      <n-space justify="end">
        <n-button
          v-if="steps.length > 0"
          type=info
          @click="nextStep()">Next Step</n-button>
        <n-button
          type=info
          @click="store.clearFront()">Skip Word</n-button>
        <n-button
          type=info
          @click="markKnown">
          Mark Known
        </n-button>
        <n-button
          type=info
          @click="submit()">Submit</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup>
import {
  toRaw, ref, reactive,
} from 'vue';
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import AnkiCardPreview from '@/components/AnkiCardPreview.vue';
import CardCreationSettings from '@/components/CardCreationSettings.vue';
import {
  useMessage, NSpace, NButton, NModal,
  NLayoutSider, NLayout, NLayoutContent,
} from 'naive-ui';
import EditSentence from '@/components/CardCreatorSteps/EditSentence.vue';
import EditImages from
  '@/components/CardCreatorSteps/EditImages.vue';
import StepsEnum from '@/components/CardCreatorSteps/StepsEnum';
import EditDefinition from
  '@/components/CardCreatorSteps/EditDefinition.vue';
import { getUserSettings } from '@/UserSettings';

const UserSettings = getUserSettings();

const store = useCardQueue();
const showModal = ref(false);
const card = ref(undefined);
let originalValues;
const step = ref(undefined);
const steps = ref([]);
let word;
let action;
let callback;
const preferBookRef = ref(undefined);
let preferBook;

// Manually change the step from an edit button.
const changeStep = (estep) => {
  step.value = estep;
  // Remove any set steps progression since the user has taken control
  steps.value = [];
};

const nextStep = async () => {
  if (steps.value.length === 0) {
    return;
  }

  const idx = steps.value.indexOf(step.value);
  if (idx === -1) {
    return;
  }
  if (idx + 1 === steps.value.length) {
    // We were on the last step
    const autoAdvanceCard = await (
      UserSettings.CardCreation.AutoAdvanceCard.read()
    );
    if (autoAdvanceCard) {
      submit();
    }
  }
  if (idx + 1 <= steps.value.length) {
    step.value = steps.value[idx + 1];
  }
};

const updateSentence = (newSentence, updateStep = false) => {
  if (newSentence.length > 0) {
    card.value.fields.sentence = newSentence;
    if (updateStep) {
      nextStep();
    }
  }
};

const updateEnglishDefinition = (newDefinitions, updateStep = false) => {
  if (newDefinitions.length > 0) {
    card.value.fields.englishDefn = newDefinitions.map(
      (def) => def.definition,
    ).join('<br>');
    const pinyin = new Set(card.value.fields.pinyin.split(', '));
    newDefinitions.forEach((def) => {
      pinyin.add(def.pronunciation);
    });
    pinyin.delete('');
    card.value.fields.pinyin = [...pinyin].join(', ');
    if (updateStep) {
      nextStep();
    }
  }
};

const updateChineseDefinition = (newDefinitions, updateStep = false) => {
  if (newDefinitions.length > 0) {
    card.value.fields.chineseDefn = newDefinitions.map(
      (def) => def.definition,
    ).join('<br>');
    const pinyin = new Set(card.value.fields.pinyin.split(', '));
    pinyin.delete('');
    newDefinitions.forEach((def) => {
      pinyin.add(def.pronunciation);
    });
    card.value.fields.pinyin = [...pinyin].join(', ');
    if (updateStep) {
      nextStep();
    }
  }
};

const updateImages = (newImages, updateStep = false) => {
  if (newImages) {
    // TODO support multiple
    card.value.fields.imageUrls = newImages.map((image) => image.thumbnailUrl);
    if (updateStep) {
      nextStep();
    }
  }
};

const message = useMessage();
store.$subscribe(async (mutation, state) => {
  // Later we can prefetch new words sentences possibly
  // if (mutation.events.type === 'add' && mutation.events.key === '0') {
  step.value = undefined;
  // TODO this is a complete mess and needs to be refined if we are going to
  // start doing anything more complicated
  if (state.wordList.length > 0) {
    [{
      word,
      action,
      preferBook,
      callback,
    }] = state.wordList;
    preferBookRef.value = preferBook;

    let ankiCard;
    if (action === ActionsEnum.CREATE) {
      ankiCard = await window.ipc.createAnkiNoteSkeleton(word);

      const enableChinese = UserSettings.Dictionaries.EnableChinese.read();
      if (enableChinese) {
        steps.value = [
          StepsEnum.SENTENCE,
          StepsEnum.ENGLISH,
          StepsEnum.CHINESE,
          StepsEnum.IMAGE,
        ];
      } else {
        steps.value = [
          StepsEnum.SENTENCE,
          StepsEnum.ENGLISH,
          StepsEnum.IMAGE,
        ];
      }
    } else {
      // Right now for EDIT we only edit the sentence so start there
      ankiCard = await window.ipc.getAnkiNote(word);
      steps.value = [
        StepsEnum.SENTENCE,
      ];
    }
    card.value = reactive(ankiCard);
    originalValues = {
      ...ankiCard.fields,
    };

    // TODO this needs to be written in a more modular way
    const englishIdx = steps.value.indexOf(StepsEnum.ENGLISH);
    if (englishIdx !== -1) {
      const autoFill = await UserSettings.CardCreation.PopulateEnglish.read();
      if (autoFill) {
        // TODO base this on default dict
        const definitions = await window.ipc.getDefinitionsForWord(
          word,
          'english',
        );
        if (definitions.length === 1) {
          updateEnglishDefinition(definitions);
          steps.value.splice(englishIdx, 1);
        }
      }
    }

    const chineseIdx = steps.value.indexOf(StepsEnum.CHINESE);
    if (chineseIdx !== -1) {
      const autoFill = await UserSettings.CardCreation.PopulateChinese.read();
      if (autoFill) {
        // TODO base this on default dict
        const definitions = await window.ipc.getDefinitionsForWord(
          word,
          'chinese',
        );
        if (definitions.length === 1) {
          updateChineseDefinition(definitions);
          steps.value.splice(chineseIdx, 1);
        }
      }
    }

    [step.value] = steps.value;
  }
  showModal.value = state.wordList.length !== 0;
});

function onClose() {
  store.clearWords();
  return false;
}

async function submit() {
  // Todo track changes to the card and submit those for update
  const messageReactive = message.create('Card submited', {
    type: 'loading',
    duration: 1e4,
  });
  // TODO figure out the logic for determining changes better
  if (action === ActionsEnum.CREATE) {
    const cardValues = toRaw(card.value.fields);
    const tags = [];
    if (preferBookRef.value !== undefined) {
      const book = await window.ipc.loadBook(preferBookRef.value);
      tags.push(book.title);
    }
    window.ipc.createAnkiCard(cardValues, tags).then((res) => {
      messageReactive.content = JSON.stringify(res);
      if (res !== 'success') {
        console.log(res);
        messageReactive.type = 'error';
      } else {
        messageReactive.type = 'success';
      }
      setTimeout(() => {
        messageReactive.destroy();
      }, 1000);
    });
  } else {
    const newData = {};
    const cardValues = toRaw(card.value.fields);
    Object.entries(cardValues).forEach(([field, value]) => {
      if (value !== originalValues[field]) {
        newData[field] = value;
      }
    });

    // Clear audio as sentence has changed
    if (newData.ExampleSentence) {
      newData.SentenceAudio = '';
    }
    console.log(toRaw(newData));
    const res = await window.ipc.updateAnkiCard(card.value.noteId, newData);
    messageReactive.content = JSON.stringify(res);
    messageReactive.type = 'success';
    setTimeout(() => {
      messageReactive.destroy();
    }, 1000);
  }
  if (callback) {
    callback();
  }
  store.clearFront();
}

function markKnown() {
  window.ipc.addWord(card.value.fields.word, 10000);
  store.clearFront();
}

</script>
