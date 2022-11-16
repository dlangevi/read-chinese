<template>
  <n-layout
    class="m-4"
    has-sider
  >
    <n-layout-sider
      width="500"
      bordered
      content-style="padding: 24px;"
    >
      <img
        class="w-auto rounded"
        :src="book.cover"
        :alt="book.title"
      >
    </n-layout-sider>
    <n-layout>
      <n-layout-header
        class="p-4"
        bordered
      >
        <p>{{ book.author }}</p>
      </n-layout-header>
      <n-layout-content class="p-8">
        <n-tabs
          type="line"
          animated
        >
          <n-tab-pane
            name="Stats"
            tab="Stats"
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
          </n-tab-pane>
          <n-tab-pane
            name="UnknownWords"
            tab="View Unknown Words"
          >
            <unknown-words
              class="h-96"
              :words="words"
              :book-filter="bookId"
            />
          </n-tab-pane>
        </n-tabs>
      </n-layout-content>
      <n-layout-footer
        class="p-4"
        bordered
      >
        <n-space justify="end">
          <button
            class="btn-primary btn"
            @click="makeFlashCards"
          >
            Make flash cards
          </button>
        </n-space>
      </n-layout-footer>
    </n-layout>
  </n-layout>
</template>

<script lang="ts" setup>
import UnknownWords from '@/components/UnknownWords.vue';
import { provide } from 'vue';
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent,
  NLayoutFooter, NSpace, NTabs, NTabPane,
} from 'naive-ui';
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

const { totalWords } = book.stats;
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
