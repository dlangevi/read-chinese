<template>
  <with-sidebar>
    <template #sidebar>
      <button class="btn-primary btn" @click="importAnki">
        Sync from Anki
      </button>
      <import-csv v-if="UserSettings.meta.EnableExperimental" />
      <select
        v-model="gridSource"
        class="select-primary select"
        @change="changeSource"
      >
        <option
          v-for="source in sources"
          :key="source"
          :value="source"
        >
          {{ source }}
        </option>
      </select>
      <div
        v-if="gridSource=='hsk'"
        class="flex place-content-between"
      >
        <span>Hsk 2.0</span>
        <input
          v-model="selectedVersion"
          type="checkbox"
          class="toggle-primary toggle toggle-lg"
        >
        <span>Hsk 3.0</span>
      </div>
      <select
        v-if="gridSource=='hsk'"
        v-model="selectedLevel"
        class="select-primary select"
        @change="loadHsk"
      >
        <option
          v-for="level in levels"
          :key="level"
          :value="level"
        >
          {{ level }}
        </option>
      </select>
      <button class="btn-primary btn" @click="makeCards">
        Make Cards
      </button>
      <div class="border-2 p-2 text-center">
        {{ words.length }} remaining words
      </div>
    </template>
    <div class="flex h-full w-full flex-col">
      <div class="shrink-0 text-center">
        <h1 class="mt-5 text-center text-xl">
          Maybe you know these words?
        </h1>
        <p>For now lets just mark learned words you already know 好不好?</p>
      </div>
      <unknown-words
        class="m-4 mx-auto h-full w-5/6 grow"
        show-definitions
        :words="words"
      />
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import WithSidebar from '@/layouts/WithSidebar.vue';
import ImportCsv from '@/components/ImportCsv.vue';
import { watch, ref, onBeforeMount } from 'vue';
import { backend } from '@wailsjs/models';
import {
  LearningTarget,
  LearningTargetFavorites,
} from '@wailsjs/backend/bookLibrary';
import {
  GetUnknownHskWords,
} from '@wailsjs/backend/KnownWords';
import { ImportAnkiKeywords } from '@wailsjs/backend/ankiInterface';
import UnknownWords from '../components/UnknownWords.vue';
import { useLoader } from '@/lib/loading';

import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import { getUserSettings } from '@/lib/userSettings';
const UserSettings = getUserSettings();

type HskVersion = '2.0' | '3.0';
type HskLevel = 1 | 2 | 3 | 4 | 5 | 6 | 7;

const version = ref<HskVersion>('2.0');
const levels = ref<HskLevel[]>([1, 2, 3, 4, 5, 6]);
const selectedLevel = ref<HskLevel>(1);
const selectedVersion = ref(false);
const active = ref(false);
watch(selectedVersion, () => {
  if (selectedVersion.value) {
    version.value = '3.0';
    levels.value = [1, 2, 3, 4, 5, 6, 7];
  } else {
    version.value = '2.0';
    levels.value = [1, 2, 3, 4, 5, 6];
  }
  if (active.value) {
    loadHsk();
  }
});

const sources = [
  'all books',
  'favorites',
  'hsk',
];

const words = ref<backend.UnknownWordEntry[]>([]);
const gridSource = ref('all books');
async function changeSource() {
  const newSource = gridSource.value;
  if (newSource === 'all books') {
    words.value = await LearningTarget();
  } else if (newSource === 'favorites') {
    words.value = await LearningTargetFavorites();
  } else if (newSource === 'hsk') {
    loadHsk();
  }
}

onBeforeMount(async () => {
  words.value = await LearningTarget();
});

const loader = useLoader();
async function importAnki() {
  return loader.withLoader(ImportAnkiKeywords, 'Syncing anki words');
}

async function loadHsk() {
  active.value = true;

  const ver = version.value;
  let lvl = selectedLevel.value;
  if (ver === '2.0' && lvl === 7) {
    lvl = 6;
  }
  words.value = await GetUnknownHskWords(ver, lvl);
}

const store = useCardQueue();
async function makeCards() {
  const sorted = words.value.slice();
  sorted.sort((a, b) => {
    if (a.occurance === undefined || b.occurance === undefined) {
      return 0;
    }
    if (a.occurance > b.occurance) {
      return -1;
    }
    return 1;
  });
  const topWords = sorted.slice(0, 50);
  topWords.forEach((word) => {
    store.addWord(word.word, ActionsEnum.CREATE);
  });
}

</script>
