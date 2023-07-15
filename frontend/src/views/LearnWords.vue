<template>
  <with-sidebar>
    <template #sidebar>
      <div class="flex max-w-full items-center gap-4">
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
      <div
        v-if="gridSource=='word list'"
        class="flex max-w-full items-center gap-4"
      >
        <label for="frequencySelect">Frequency List:</label>
        <select
          id="frequencySelect"
          v-model="frequencyList"
          class="select-primary select w-1/2 grow"
          @change="changeSource"
        >
          <option
            v-for="source in frequencyLists"
            :key="source"
            :value="source"
            class="max-w-full"
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
      <div
        v-if="gridSource=='hsk'"
        class="border-2 p-2 text-center"
      >
        {{ words.length }} remaining words
      </div>
      <div class="flex items-center gap-4">
        <button
          class="btn-primary btn w-1/3"
          @click="changeSort(false)"
        >
          Sort by occurance
        </button>
        <button
          class="btn-primary btn w-1/3"
          @click="changeSort(true)"
        >
          Sort by frequency
        </button>
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
        :sort-by-frequency="sortByFrequency"
        :frequency-source="frequencySource"
        :occurance-source="occuranceSource"
      />
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import WithSidebar from '@/layouts/WithSidebar.vue';
import { watch, ref, onBeforeMount } from 'vue';
import {
  LearningTarget,
  LearningTargetFavorites,
} from '@wailsjs/backend/bookLibrary';
import {
  GetUnknownHskWords,
  GetUnknownListWords,
  GetLists,
  GetPrimaryList,
} from '@wailsjs/backend/wordLists';
import UnknownWords from '../components/UnknownWords.vue';
import { GetPossibleWords } from '@wailsjs/backend/Dictionaries';

type HskVersion = '2.0' | '3.0';
type HskLevel = 1 | 2 | 3 | 4 | 5 | 6 | 7;

const version = ref<HskVersion>('2.0');
const levels = ref<HskLevel[]>([1, 2, 3, 4, 5, 6]);
const selectedLevel = ref<HskLevel>(1);
const selectedVersion = ref(false);
const active = ref(false);
const unknownWordsRef = ref();
const sortByFrequency = ref(false);
function changeSort(freq : boolean) {
  console.log('changing sort', freq);
  sortByFrequency.value = freq;
}

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
  'word list',
];

const words = ref<string[]>([]);
const gridSource = ref('all books');
const frequencyLists = ref<string[]>([]);
let primaryFrequency = '';
// The chosen frequencyList to mine
const frequencyList = ref('');
// The source of frequency lists
const frequencySource = ref('');
const occuranceSource = ref('all');
const searchBox = ref('');
async function changeSource() {
  frequencySource.value = primaryFrequency;
  const newSource = gridSource.value;
  if (newSource === 'all books') {
    words.value = await LearningTarget();
    occuranceSource.value = 'all';
  } else if (newSource === 'favorites') {
    words.value = await LearningTargetFavorites();
    occuranceSource.value = 'favorites';
  } else if (newSource === 'hsk') {
    loadHsk();
    occuranceSource.value = 'all';
  } else if (newSource === 'word list') {
    loadWordList();
    occuranceSource.value = 'all';
  }
}

onBeforeMount(async () => {
  words.value = await LearningTarget();
  frequencyLists.value = await GetLists();
  primaryFrequency = await GetPrimaryList();
  frequencyList.value = primaryFrequency;
  frequencySource.value = primaryFrequency;
  occuranceSource.value = 'all';
});

async function loadWordList() {
  // Set the frequency source here
  frequencySource.value = frequencyList.value;
  words.value = await GetUnknownListWords(frequencyList.value);
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
