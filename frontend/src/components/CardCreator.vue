<template>
  <div
    :class="['modal', {'modal-open': showModal}]"
    @click="onClose"
  >
    <div
      v-if="showModal"
      class="modal-box flex h-[80vh] w-4/5 max-w-full flex-col"
      @click.stop
    >
      <!-- TODO dont use absolute positions. Do like a flex on the right -->
      <div class="absolute right-4 top-4 flex gap-2">
        <card-creation-settings />
        <button
          class="btn-sm btn-circle btn"
          @click="onClose"
        >
          ✕
        </button>
      </div>
      <p class="text-xl">
        Creating card for {{ cardManager.word }}
      </p>
      <div class="flex">
        <!-- v-if="card !== undefined" -->
        <div class="w-1/3">
          <anki-card-preview @change-step="changeStep" />
        </div>
        <div class="h-[60vh] w-full overflow-scroll p-4">
          <p class="text-4xl">
            {{ cardManager.word }}
          </p>
          <edit-sentence
            v-show="step === StepsEnum.SENTENCE"
            v-if="steps.includes(StepsEnum.SENTENCE)"
            :prefer-book="preferBookRef"
          />
          <edit-definition
            v-show="step === StepsEnum.ENGLISH"
            v-if="steps.includes(StepsEnum.ENGLISH)"
            type="english"
          />
          <edit-definition
            v-show="step === StepsEnum.CHINESE"
            v-if="steps.includes(StepsEnum.CHINESE)"
            type="chinese"
          />
          <edit-images
            v-show="step === StepsEnum.IMAGE"
            v-if="steps.includes(StepsEnum.IMAGE)"
          />
        </div>
      </div>

      <div class="modal-action">
        <div class="flex place-content-end gap-2">
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
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  ref, toRaw,
} from 'vue';
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import { useCardManager } from '@/stores/CardManager';
import AnkiCardPreview from '@/components/AnkiCardPreview.vue';
import CardCreationSettings from '@/components/CardCreationSettings.vue';
import { useMessage } from '@/lib/messages';
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
} from '@wailsjs/backend/ankiInterface';

import {
  GetBook,
} from '@wailsjs/backend/bookLibrary';
import {
  backend,
} from '@wailsjs/models';

const UserSettings = getUserSettings();

const store = useCardQueue();
const cardManager = useCardManager();
const showModal = ref(false);
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
  // Dont want to auto advance if we are modifying a card
  if (action === ActionsEnum.MODIFY) {
    step.value = steps.value[idx + 1];
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

// TODO I am leaving the commented out old message code with the plan
// of eventually having that sort of api avaliable
const message = useMessage();
store.$subscribe(async (_, state) => {
  // This is needed to reset the state
  step.value = StepsEnum.NONE;
  // Also needed to reset the state since it flips the v-if statements
  // TODO this is hacky and needs a better solution
  steps.value = [];
  // Later we can prefetch new words sentences possibly
  // if (mutation.events.type === 'add' && mutation.events.key === '0') {
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
    } else {
      // Right now for EDIT we only edit the sentence so start there
      ankiCard = await GetAnkiNote(word);
    }
    cardManager.loadCard(ankiCard);
    // TODO figure out how to weave step tracking between
    // the two. One idea would be to watch for getters to become set
    steps.value = cardManager.steps;

    stepsFilled = { };
    steps.value.forEach(step => {
      stepsFilled[step] = false;
    });
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
  // const messageReactive = message.create('Card submited', {
  //   type: 'loading',
  //   duration: 1e4,
  // });
  message.info('Card submited');
  // TODO figure out the logic for determining changes better
  if (action === ActionsEnum.CREATE) {
    const cardValues = toRaw(cardManager.newValues);
    const tags = [];
    if (preferBookRef.value !== undefined) {
      const book = await GetBook(preferBookRef.value);
      tags.push(book.title);
    }
    // TODO do this kind of catching elsewhere
    console.log('creating values ', cardValues);
    CreateAnkiNote(cardValues, tags)
      .then(() => {
        message.success('success');
      })
      .catch((err) => {
        message.error(err);
      });
  } else {
    const cardValues: backend.Fields = toRaw(cardManager.newValues);
    UpdateNoteFields(cardManager.note.noteId, cardValues)
      .then(() => {
        message.success('success');
      })
      .catch((err) => {
        message.error(err);
      });
  }
  if (callback) {
    callback();
  }
  store.clearFront();
}

function markKnown() {
  AddWord(cardManager.word, 10000);
  store.clearFront();
}
</script>
