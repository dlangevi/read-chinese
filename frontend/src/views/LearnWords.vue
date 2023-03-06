<template>
  <with-sidebar>
    <template #sidebar>
      <div class="flex items-center gap-4">
        <label for="sourceSelect">Source:</label>
        <select
          id="sourceSelect"
          v-model="gridSource"
          class="select-primary select grow"
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
      </div>
      <div v-if="gridSource=='search'">
        <input
          v-model="searchBox"
          type="text"
          placeholder="Input here (using 汉字 for now)"
          class="input-primary input mx-auto w-full"
          @input="onUpdateSearchBox"
        >
      </div>
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
      <button class="btn-primary btn" @click="importAnki">
        Sync from Anki
      </button>
      <import-csv v-if="UserSettings.meta.EnableExperimental" />
      <div
        v-if="gridSource=='hsk'"
        class="border-2 p-2 text-center"
      >
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
        ref="unknownWordsRef"
        class="m-4 mx-auto h-full w-5/6 grow"
        show-definitions
        :words="words"
        :frequency-source="frequencySource"
      />
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import WithSidebar from '@/layouts/WithSidebar.vue';
import ImportCsv from '@/components/ImportCsv.vue';
import { watch, ref, onBeforeMount } from 'vue';
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
import { GetPossibleWords } from '@wailsjs/backend/Dictionaries';

import { getUserSettings } from '@/lib/userSettings';
const UserSettings = getUserSettings();

type HskVersion = '2.0' | '3.0';
type HskLevel = 1 | 2 | 3 | 4 | 5 | 6 | 7;

const version = ref<HskVersion>('2.0');
const levels = ref<HskLevel[]>([1, 2, 3, 4, 5, 6]);
const selectedLevel = ref<HskLevel>(1);
const selectedVersion = ref(false);
const active = ref(false);
const unknownWordsRef = ref();

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
  'search',
];

const words = ref<string[]>([]);
const gridSource = ref('all books');
const frequencySource = ref('');
const searchBox = ref('');
async function changeSource() {
  const newSource = gridSource.value;
  if (newSource === 'all books') {
    words.value = await LearningTarget();
    frequencySource.value = '';
  } else if (newSource === 'favorites') {
    words.value = await LearningTargetFavorites();
    frequencySource.value = 'favorites';
  } else if (newSource === 'hsk') {
    loadHsk();
    frequencySource.value = '';
  }
}

onBeforeMount(async () => {
  words.value = await LearningTarget();
  frequencySource.value = '';
});

const loader = useLoader();
async function importAnki() {
  return loader.withLoader(ImportAnkiKeywords);
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

function makeCards() {
  unknownWordsRef.value.enqueueTopRows(50);
}

async function onUpdateSearchBox() {
  if (searchBox.value.length !== 0) {
    words.value = await GetPossibleWords(searchBox.value);
  }
}

</script>
