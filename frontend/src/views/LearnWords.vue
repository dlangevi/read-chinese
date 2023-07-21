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
      <div v-if="unknownWordsRef" class="border-2 p-2 text-center">
        {{ unknownWordsRef.getRowCount() }} remaining words
      </div>
      <div class="flex items-center gap-4">
        <button
          class="btn-primary btn w-1/3"
          @click="SortByOccurance"
        >
          Sort by occurance
        </button>
        <button
          class="btn-primary btn w-1/3"
          @click="SortByFrequency"
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
      />
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import WithSidebar from '@/layouts/WithSidebar.vue';
import { watch, ref, onBeforeMount } from 'vue';
import {
  GetLists,
  GetPrimaryList,
  SetFrequencySource,
  SetOccuranceSource,
  SetWordSourceFromAll,
  SetWordSourceFromFavorites,
  SetWordSourceFromHsk,
  SetWordSourceFromSearch,
  SetWordSourceFromList,
  SortByFrequency,
  SortByOccurance,
} from '@wailsjs/backend/wordLists';
import UnknownWords from '../components/UnknownWords.vue';

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
  'word list',
];

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
    SetWordSourceFromAll();
  } else if (newSource === 'favorites') {
    SetWordSourceFromFavorites();
  } else if (newSource === 'hsk') {
    loadHsk();
  } else if (newSource === 'word list') {
    frequencySource.value = frequencyList.value;
    SetFrequencySource(frequencyList.value);
    SetWordSourceFromList(frequencyList.value);
  }
}

onBeforeMount(async () => {
  frequencyLists.value = await GetLists();
  primaryFrequency = await GetPrimaryList();
  frequencyList.value = primaryFrequency;
  frequencySource.value = primaryFrequency;
  occuranceSource.value = 'all';
  await SetOccuranceSource('all');
  await SetFrequencySource(primaryFrequency);
  changeSource();
});

async function loadHsk() {
  active.value = true;
  const ver = version.value;
  let lvl = selectedLevel.value;
  if (ver === '2.0' && lvl === 7) {
    lvl = 6;
  }
  SetWordSourceFromHsk(ver, lvl);
}

function makeCards() {
  unknownWordsRef.value.enqueueTopRows(500);
}

async function onUpdateSearchBox() {
  if (searchBox.value.length !== 0) {
    SetWordSourceFromSearch(searchBox.value);
  }
}

</script>
