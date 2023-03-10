<template>
  <div
    class="flex"
    has-sider
  >
    <div
      class="basis-1/4"
      bordered
      content-style="padding: 24px;"
    >
      <img
        class="w-auto rounded"
        :src="'/' + book.cover"
        :alt="book.title"
      >
    </div>
    <div class="grow basis-3/4">
      <div
        class="p-4"
        bordered
      >
        <p class="text-3xl">{{ book.author }} - {{ book.title }}</p>
      </div>
      <div class="p-8">
        <tabbed-pane>
          <tabbed-pane-tab
            name="Stats"
            title="Stats"
          >
            <div class="stats shadow">
              <div class="stat place-items-center">
                <div class="stat-title"> Known </div>
                <div class="stat-value"> {{ known.toFixed(2) }}% </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Can Read </div>
                <div class="stat-value"> {{ likelyKnown.toFixed(2) }}% </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Known Characters </div>
                <div class="stat-value">
                  {{ knownCharacters.toFixed(2) }}%
                </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Total Words </div>
                <div class="stat-value"> {{ stats?.totalWords }} </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Total Characters </div>
                <div class="stat-value"> {{ stats?.totalCharacters }} </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Unique Characters </div>
                <div class="stat-value"> {{ stats?.uniqueCharacters }} </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Unique Words </div>
                <div class="stat-value"> {{ stats?.uniqueWords }} </div>
              </div>
            </div>

            <table class="table w-full">
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
          </tabbed-pane-tab>
          <tabbed-pane-tab
            name="UnknownWords"
            title="View Unknown Words"
          >
            <unknown-words
              class="h-96"
              :words="words"
              :book-filter="bookId"
            />
          </tabbed-pane-tab>
        </tabbed-pane>
      </div>
      <div
        class="p-4"
        bordered
      >
        <div class="flex place-content-end">
          <button
            class="btn-primary btn"
            @click="makeFlashCards"
          >
            Make flash cards
          </button>
        </div>
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
  onUnmounted, onBeforeMount, ref, computed,
} from 'vue';

import {
  TopUnknownWords,
  LearningTargetBook,
  GetBook,
} from '@wailsjs/backend/bookLibrary';

const props = defineProps({
  bookId: {
    required: true,
    type: Number,
  },
});

const words = ref<string[]>([]);
const book = ref<backend.Book>(backend.Book.createFrom());
const stats = computed(() => book.value.stats);

const known = ref(0);
const likelyKnown = ref(0);
const knownCharacters = ref(0);
const targetPairs = ref<{
  target: number,
  words: number,
}[]>([]);

async function loadBook() {
  book.value = await GetBook(props.bookId);
  const { stats } = book.value;

  known.value = (
    (stats.totalKnownWords / stats.totalWords) * 100);
  likelyKnown.value = (
    (stats.probablyKnownWords / stats.totalWords) * 100);
  knownCharacters.value = (
    (stats.knownCharacters / stats.totalCharacters) * 100);

  const firstTarget = stats.needToKnow.findIndex((n) => n !== 0);

  const targets = stats.targets.slice(
    firstTarget, firstTarget + 3);

  const needToKnow = stats.needToKnow.slice(
    firstTarget, firstTarget + 3);

  targetPairs.value = targets.map((e, i) => (
    { target: e, words: needToKnow[i] }));
}

onBeforeMount(async () => {
  await loadBook();
  words.value = await LearningTargetBook(book.value.bookId);
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
