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
                <div class="stat-value"> {{ known }}% </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Can Read </div>
                <div class="stat-value"> {{ likelyKnown }}% </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Known Characters </div>
                <div class="stat-value"> {{ knownCharacters }}% </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Total Words </div>
                <div class="stat-value"> {{ totalWords }} </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Total Characters </div>
                <div class="stat-value"> {{ totalCharacters }} </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Unique Characters </div>
                <div class="stat-value"> {{ uniqueCharacters }} </div>
              </div>
              <div class="stat place-items-center">
                <div class="stat-title"> Unique Words </div>
                <div class="stat-value"> {{ uniqueWords }} </div>
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
import { provide } from 'vue';
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import type { backend } from '@wailsjs/models';

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

provide('preferBook', props.bookId);

const book:backend.Book = await GetBook(props.bookId);
console.log(book);
const words = await LearningTargetBook(book.bookId);
const { stats } = book;

const known = (
  (stats.totalKnownWords / stats.totalWords) * 100).toFixed(2);
const likelyKnown = (
  (stats.probablyKnownWords / stats.totalWords) * 100).toFixed(2);
const knownCharacters = (
  (stats.knownCharacters / stats.totalCharacters) * 100).toFixed(2);

const {
  totalWords, totalCharacters,
  uniqueWords, uniqueCharacters,
} = book.stats;
const firstTarget = stats.needToKnow.findIndex((n) => n !== 0);
const targets = stats.targets.slice(firstTarget, firstTarget + 3);
const needToKnow = stats.needToKnow.slice(firstTarget, firstTarget + 3);
const targetPairs = targets.map((e, i) => (
  { target: e, words: needToKnow[i] }));

const store = useCardQueue();
async function makeFlashCards() {
  const topWords: string[] = await TopUnknownWords(props.bookId, 50);
  topWords.forEach((word) => {
    store.addWord(word, ActionsEnum.CREATE, { preferBook: props.bookId });
  });
}
</script>
