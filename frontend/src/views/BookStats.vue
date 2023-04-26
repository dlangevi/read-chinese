<template>
  <div
    ref="styleRef"
    class="flex"
  >
    <div class="h-full w-1/4 p-4">
      <img
        v-if="book.cover"
        class="w-auto rounded"
        :src="'/' + book.cover"
        :alt="book.title"
      >
      <img
        v-else
        class="w-auto rounded"
        src="../assets/empty-book.jpg"
        alt="Default image"
      >
    </div>
    <div class="flex h-full w-3/4 flex-col">
      <div
        class="flex place-items-center p-4"
        bordered
      >
        <h2 class="text-3xl">{{ book.author }} - {{ book.title }}</h2>
        <div class="flex grow place-content-end">
          <button
            class="btn-primary btn"
            @click="makeFlashCards"
          >
            Make flash cards
          </button>
        </div>
      </div>
      <div class="grow p-8">
        <tabbed-pane class="h-full">
          <tabbed-pane-tab
            name="Stats"
            tab-title="Stats"
          >
            <div class="flex h-full gap-4">
              <div class="flex w-1/2 flex-col gap-4">
                <div class="stats shadow">
                  <div
                    v-for="val, title in {
                      'Total Words': book.totalWords,
                      'Total Characters': book.totalCharacters,
                      'Unique Characters': book.uniqueCharacters,
                      'Unique Words': book.uniqueWords,
                    }"
                    :key="title"
                    class="stat place-items-center"
                  >
                    <div class="stat-title"> {{ title }} </div>
                    <div class="stat-value"> {{ val }} </div>
                  </div>
                </div>

                <div class="stats shadow">
                  <div
                    v-for="val, title in {
                      Known: known.toFixed(2),
                      'Can Read': likelyKnown.toFixed(2),
                      'Known Characters': knownCharacters.toFixed(2),
                    }"
                    :key="title"
                    class="stat place-items-center"
                  >
                    <div class="stat-title"> {{ title }} </div>
                    <div class="stat-value"> {{ val }}% </div>
                  </div>
                </div>

                <table class="table">
                  <thead>
                    <tr>
                      <th>Target</th>
                      <th>Needed Words</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="pair in targetPairs"
                      :key="pair.target"
                      :label="pair.target"
                      :value="pair.words"
                    >
                      <td>{{ pair.target }}</td>
                      <td>{{ pair.words }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <div class="h-3/4 w-1/2">
                <Line
                  v-if="loaded"
                  :data="data"
                  :options="options"
                />
              </div>
            </div>
          </tabbed-pane-tab>
          <tabbed-pane-tab
            name="UnknownWords"
            tab-title="View Unknown Words"
          >
            <unknown-words
              class="h-4/5"
              :words="words"
              :book-filter="bookId"
            />
          </tabbed-pane-tab>
        </tabbed-pane>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import TabbedPane from '@/layouts/TabbedPane.vue';
import TabbedPaneTab from '@/components/TabbedPaneTab.vue';
import UnknownWords from '@/components/UnknownWords.vue';
import { useCardQueue } from '@/stores/CardQueue';
import { backend } from '@wailsjs/models';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import {
  onUnmounted, onMounted, ref,
} from 'vue';

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import type { ChartData } from 'chart.js';
import { Line } from 'vue-chartjs';

import {
  TopUnknownWords,
  LearningTargetBook,
  GetBookGraph,
  GetBook,
} from '@wailsjs/backend/bookLibrary';

// TODO this is jank and global
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
);

const options = {
  responsive: true,
  maintainAspectRatio: false,
};

const props = defineProps({
  bookId: {
    required: true,
    type: Number,
  },
});

const words = ref<string[]>([]);
const book = ref<backend.Book>(backend.Book.createFrom());
const known = ref(0);
const likelyKnown = ref(0);
const knownCharacters = ref(0);
const targetPairs = ref<{
  target: number,
  words: number,
}[]>([]);
const data = ref<ChartData<'line'>>({
  datasets: [],
});
const loaded = ref(false);
const styleRef = ref<HTMLElement | null>(null);

async function loadBook() {
  book.value = await GetBook(props.bookId);
  const { stats } = book.value;

  known.value = (
    (stats.totalKnownWords / book.value.totalWords) * 100);
  likelyKnown.value = (
    (stats.probablyKnownWords / book.value.totalWords) * 100);
  knownCharacters.value = (
    (stats.knownCharacters / book.value.totalCharacters) * 100);

  const firstTarget = stats.needToKnow.findIndex((n) => n !== 0);

  const targets = stats.targets.slice(
    firstTarget, firstTarget + 3);

  const needToKnow = stats.needToKnow.slice(
    firstTarget, firstTarget + 3);

  targetPairs.value = targets.map((e, i) => (
    { target: e, words: needToKnow[i] }));
}

onMounted(async () => {
  if (!styleRef.value) {
    console.error('styleref not there');
    return;
  }
  await loadBook();
  words.value = await LearningTargetBook(book.value.bookId);
  const rawdata = await GetBookGraph(book.value.bookId);
  const style = getComputedStyle(styleRef.value);
  const primary = `hsl(${style.getPropertyValue('--p')}`;
  data.value = {
    labels: rawdata.map(d => d.day),
    datasets: [{
      label: 'Percent known over time',
      backgroundColor: primary,
      borderColor: primary,
      pointRadius: 0,
      tension: 0.5,
      data: rawdata.map(d => d.known),
    }],
  };
  loaded.value = true;

  EventsOn('AddedWord', async () => {
    loadBook();
  });
});

onUnmounted(() => {
  EventsOff('AddedWord');
});

const store = useCardQueue();
async function makeFlashCards() {
  const topWords: string[] = await TopUnknownWords(props.bookId, 50);
  topWords.forEach((word) => {
    store.addWord({ word }, () => {}, props.bookId);
  });
}
</script>
