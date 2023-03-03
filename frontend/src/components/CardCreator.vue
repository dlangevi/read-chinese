<template>
  <div
    :class="['modal', {'modal-open': showModal}]"
    @click.stop
  >
    <div
      v-if="showModal"
      class="modal-box flex h-[80vh] w-4/5 max-w-full flex-col"
      @click.stop
    >
      <h2 class="m-4 text-4xl font-semibold">
        Creating card for {{ cardManager.word }}
      </h2>
      <div class="flex">
        <anki-card-preview class="w-1/3" />
        <div
          :key="cardManager.word"
          class="h-[60vh] w-full overflow-scroll p-4"
        >
          <edit-sentence
            v-show="cardManager.currentStep === StepsEnum.SENTENCE"
            v-if="cardManager.steps.includes(StepsEnum.SENTENCE)"
            :prefer-book="preferBook"
          />
          <edit-definition
            v-show="cardManager.currentStep === StepsEnum.ENGLISH"
            v-if="cardManager.steps.includes(StepsEnum.ENGLISH)"
            :english="true"
          />
          <edit-definition
            v-show="cardManager.currentStep === StepsEnum.CHINESE"
            v-if="cardManager.steps.includes(StepsEnum.CHINESE)"
            :chinese="true"
          />
          <edit-images
            v-show="cardManager.currentStep === StepsEnum.IMAGE"
            v-if="cardManager.steps.includes(StepsEnum.IMAGE)"
          />
        </div>
      </div>

      <!-- Settings and close buttons -->
      <div class="absolute right-4 top-4 flex gap-2">
        <card-creation-settings />
        <close-circle-sharp class="h-6 w-6 cursor-pointer" @click="onClose" />
      </div>

      <div class="modal-action flex items-center justify-center">
        <div class="w-1/4" />
        <card-creator-tracker class="grow justify-start" />
        <div class="flex place-content-end gap-2">
          <button
            class="btn-primary btn-sm btn"
            @click="store.clearFront"
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
            @click="submit"
          >
            Submit
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { CloseCircleSharp } from '@vicons/ionicons5';
import {
  ref, watch,
} from 'vue';
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import { useCardManager } from '@/stores/CardManager';
import AnkiCardPreview from '@/components/AnkiCardPreview.vue';
import CardCreationSettings from '@/components/CardCreationSettings.vue';
import { useMessage } from '@/lib/messages';
import CardCreatorTracker from
  '@/components/CardCreatorSteps/CardCreatorTracker.vue';
import EditSentence from '@/components/CardCreatorSteps/EditSentence.vue';
import EditImages from
  '@/components/CardCreatorSteps/EditImages.vue';
import EditDefinition from
  '@/components/CardCreatorSteps/EditDefinition.vue';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';

import { AddWord } from '@wailsjs/backend/KnownWords';
import {
  CreateAnkiNote,
  UpdateNoteFields,
  GetAnkiNote,
  GetAnkiNoteSkeleton,
} from '@wailsjs/backend/ankiInterface';

import {
  UpdateSentenceTable,
} from '@wailsjs/backend/Generator';
import {
  backend,
} from '@wailsjs/models';

import { storeToRefs } from 'pinia';
import { getUserSettings } from '@/lib/userSettings';

const store = useCardQueue();
const cardManager = useCardManager();
const showModal = ref(false);
const action = ref<ActionsEnum>(ActionsEnum.CREATE);
const callback = ref<(() => void) | undefined>();
const preferBook = ref<number | undefined>();

const UserSettings = getUserSettings();
const message = useMessage();
store.$subscribe(async (_, state) => {
  if (state.wordList.length > 0) {
    const [{
      word,
      action: _action,
      preferBook: _preferBook,
      callback: _callback,
    }] = state.wordList;
    action.value = _action;
    callback.value = _callback;
    preferBook.value = _preferBook;

    let ankiCard;
    if (action.value === ActionsEnum.CREATE) {
      ankiCard = await GetAnkiNoteSkeleton(word);
    } else {
      ankiCard = await GetAnkiNote(word);
    }
    const enableChinese = UserSettings.Dictionaries.EnableChinese;
    const hasImageApi = UserSettings.AnkiConfig.AzureImageApiKey !== '';
    cardManager.loadCard(ankiCard, enableChinese, hasImageApi);
  }
  showModal.value = state.wordList.length !== 0;
});

function onClose() {
  store.clearWords();
  return false;
}

const { ready } = storeToRefs(cardManager);
watch(ready, () => {
  const autoSubmit = UserSettings.CardCreation.AutoAdvanceCard;
  if (ready.value && autoSubmit) {
    submit();
  }
});

async function submit() {
  message.info('Card submited');
  if (action.value === ActionsEnum.CREATE) {
    const cardValues = cardManager.getChanged();

    const tags : string[] = [];
    if (UserSettings.AnkiConfig.AddBookTag) {
      if (cardManager.sentenceSource) {
        tags.push(cardManager.sentenceSource);
      }
    }
    if (UserSettings.AnkiConfig.AddProgramTag) {
      tags.push('read-chinese');
    }

    console.log('creating values ', cardValues);
    await CreateAnkiNote(cardValues, tags);
    message.success('success');
  } else {
    const cardValues: backend.Fields = cardManager.getChanged();
    await UpdateNoteFields(cardManager.note.noteId, cardValues);
    message.success('success');
  }
  if (callback.value) {
    callback.value();
  }
  store.clearFront();
}

async function markKnown() {
  await AddWord(cardManager.word, 10000);
  UpdateSentenceTable(cardManager.word);
  store.clearFront();
}
</script>
