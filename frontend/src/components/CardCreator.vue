<template>
  <n-modal
    v-model:show="showModal"
    class="h-[80vh] w-4/5"
    :mask-closable="false"
    :closable="true"
    preset="card"
    :on-close="onClose"
    content-style="display:flex; flex-direction: column;"
  >
    <template #header>
      <p class="text-xl">
        Creating card for {{ card.fields.word }}
      </p>
    </template>
    <template #header-extra>
      <card-creation-settings />
    </template>
    <n-layout
      has-sider
      sider-placement="left"
      class="grow"
      style="max-height: 60vh"
    >
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
          <anki-card-preview
            :anki-card="card"
            @change-step="changeStep"
          />
        </Suspense>
      </n-layout-sider>
      <n-layout-content
        content-style="padding: 24px;"
        :native-scrollbar="true"
      >
        <p class="text-4xl">
          {{ card.fields.word }}
        </p>
        <edit-sentence
          v-show="step === StepsEnum.SENTENCE"
          v-if="steps.includes(StepsEnum.SENTENCE)"
          :prefer-book="preferBookRef"
          :word="card.fields.word"
          :sentence="card.fields.sentence"
          @update-sentence="updateSentence"
        />
        <edit-definition
          v-show="step === StepsEnum.ENGLISH"
          v-if="steps.includes(StepsEnum.ENGLISH)"
          :word="card.fields.word"
          :definition="card.fields.englishDefn"
          type="english"
          @update-definition="updateEnglishDefinition"
        />
        <edit-definition
          v-show="step === StepsEnum.CHINESE"
          v-if="steps.includes(StepsEnum.CHINESE)"
          :word="card.fields.word"
          :definition="card.fields.englishDefn"
          type="chinese"
          @update-definition="updateChineseDefinition"
        />
        <edit-images
          v-show="step === StepsEnum.IMAGE"
          v-if="steps.includes(StepsEnum.IMAGE)"
          :word="card.fields.word"
          @update-images="updateImages"
        />
      </n-layout-content>
    </n-layout>

    <template #action>
      <n-space justify="end">
        <button
          v-if="steps.length > 0"
          class="btn-primary btn-sm btn"
          @click="nextStep()"
        >
          Next Step
        </button>
        <button
          class="btn-primary btn-sm btn"
          @click="store.clearFront()"
        >
          Skip Word
        </button>
        <button
          class="btn-primary btn-sm btn"
          @click="markKnown"
        >
          Mark Known
        </button>
        <button
          class="btn-primary btn-sm btn"
          @click="submit()"
        >
          Submit
        </button>
      </n-space>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
import {
  toRaw, ref, reactive,
} from 'vue';
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import AnkiCardPreview from '@/components/AnkiCardPreview.vue';
import CardCreationSettings from '@/components/CardCreationSettings.vue';
import {
  useMessage, NSpace, NModal,
  NLayoutSider, NLayout, NLayoutContent,
} from 'naive-ui';
import EditSentence from '@/components/CardCreatorSteps/EditSentence.vue';
import EditImages from
  '@/components/CardCreatorSteps/EditImages.vue';
import EditDefinition from
  '@/components/CardCreatorSteps/EditDefinition.vue';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import { getUserSettings } from '@/lib/userSettings';

import { AddWord } from '@wailsjs/backend/KnownWords';
import {
  CreateAnkiNote,
  UpdateNoteFields,
  GetAnkiNote,
  GetAnkiNoteSkeleton,
} from '@wailsjs/backend/AnkiInterface';

import {
  GetBook,
} from '@wailsjs/backend/bookLibrary';
import {
  backend,
} from '@wailsjs/models';

const UserSettings = getUserSettings();

const store = useCardQueue();
const showModal = ref(false);
const card = ref<backend.RawAnkiNote>(new backend.RawAnkiNote());
const step = ref<StepsEnum>(StepsEnum.SENTENCE);
const steps = ref<StepsEnum[]>([]);
let stepsFilled : { [s:string]: boolean } = {};
let word;
let action: ActionsEnum;
let callback: (() => void) | undefined;
const preferBookRef = ref<number | undefined>(undefined);
let preferBook;

// Manually change the step from an edit button.
const changeStep = (estep: StepsEnum) => {
  step.value = estep;
};

const nextStep = async () => {
  const idx = steps.value.indexOf(step.value);
  if (idx === -1) {
    return;
  }
  if (idx + 1 === steps.value.length) {
    // We were on the last step
    if (UserSettings.CardCreation.AutoAdvanceCard.read()) {
      submit();
    }
  }
  step.value = steps.value[idx + 1];
  if (stepsFilled[step.value]) {
    nextStep();
  }
};

const updateSentence = (newSentence: string, updateStep = false) => {
  if (newSentence.length > 0) {
    card.value.fields.sentence = newSentence;
    stepsFilled[StepsEnum.SENTENCE] = true;
    if (updateStep) {
      nextStep();
    }
  }
};

const updateDefinition = (
  newDefinitions: backend.DictionaryEntry[],
  updateStep: boolean,
  setter: (arg0: string) => void,
) => {
  const definitions = newDefinitions.map(
    (def) => `[${def.pronunciation}] ${def.definition}`,
  ).join('<br>');
  setter(definitions);
  const pinyin = new Set(card.value.fields.pinyin.split(', '));
  newDefinitions.forEach((def) => {
    const pronunciation = def.pronunciation.replace(/\s/g, '');
    pinyin.add(pronunciation);
  });
  pinyin.delete('');
  card.value.fields.pinyin = [...pinyin].join(', ');
  if (updateStep) {
    nextStep();
  }
};

const updateEnglishDefinition = (
  newDefinitions: backend.DictionaryEntry[],
  updateStep = false,
) => {
  updateDefinition(newDefinitions, updateStep, function (newDefs: string) {
    card.value.fields.englishDefn = newDefs;
  });
  stepsFilled[StepsEnum.ENGLISH] = true;
};

const updateChineseDefinition = (
  newDefinitions: backend.DictionaryEntry[],
  updateStep = false,
) => {
  updateDefinition(newDefinitions, updateStep, function (newDefs: string) {
    card.value.fields.chineseDefn = newDefs;
  });
  stepsFilled[StepsEnum.CHINESE] = true;
};

const updateImages = (newImages: backend.ImageInfo[], updateStep = false) => {
  if (newImages) {
    card.value.fields.imageUrls = newImages.map((image) => image.thumbnailUrl);
    if (updateStep) {
      nextStep();
    }
  }
  stepsFilled[StepsEnum.IMAGE] = true;
};

const message = useMessage();
store.$subscribe(async (_, state) => {
  // This is needed to reset the state
  step.value = StepsEnum.NONE;
  // Also needed to reset the state since it flips the v-if statements
  // TODO this is hacky and needs a better solution
  steps.value = [];
  // Later we can prefetch new words sentences possibly
  // if (mutation.events.type === 'add' && mutation.events.key === '0') {
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
      ankiCard = await GetAnkiNoteSkeleton(word);

      // Todo base this on default dict
      steps.value = [
        StepsEnum.SENTENCE,
        StepsEnum.ENGLISH,
        ...(
          // Lol javascript is fun
          UserSettings.Dictionaries.EnableChinese.read()
            ? [StepsEnum.CHINESE]
            : []
        ),
        StepsEnum.IMAGE,
      ];
    } else {
      // Right now for EDIT we only edit the sentence so start there
      ankiCard = await GetAnkiNote(word);
      steps.value = [
        StepsEnum.SENTENCE,
      ];
    }
    stepsFilled = { };
    steps.value.forEach(step => {
      stepsFilled[step] = false;
    });
    card.value = reactive(ankiCard);
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
      const book = await GetBook(preferBookRef.value);
      tags.push(book.title);
    }
    // TODO do this kind of catching elsewhere
    CreateAnkiNote(cardValues, tags)
      .then(() => {
        messageReactive.content = 'success';
        messageReactive.type = 'success';
        setTimeout(() => {
          messageReactive.destroy();
        }, 1000);
      })
      .catch((err) => {
        messageReactive.content = err;
        messageReactive.type = 'error';
        setTimeout(() => {
          messageReactive.destroy();
        }, 1000);
      });
  } else {
    const cardValues: backend.Fields = toRaw(card.value.fields);
    const res = await UpdateNoteFields(card.value.noteId, cardValues);
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
  AddWord(card.value.fields.word, 10000);
  store.clearFront();
}
</script>
