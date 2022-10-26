<template>
  <n-layout class="m-4" has-sider>
    <n-layout-sider width=500 bordered content-style="padding: 24px;">
      <img
        class="rounded rounded-t w-auto"
        :src="assetUrl"
        :alt="book.title"
      />
    </n-layout-sider>
    <n-layout>
      <n-layout-header class="p-4" bordered>
        <p>{{ book.author }}</p>
      </n-layout-header>
      <n-layout-content class="p-8">
        <n-tabs type="line" animated>
          <n-tab-pane name="Stats" tab="Stats">
            <n-statistic label="Known" :value="known">
              <template #suffix>
                %
              </template>
            </n-statistic>
            <n-statistic label="Can Read" :value="likelyKnown">
              <template #suffix>
                %
              </template>
            </n-statistic>
            <n-statistic label="Known Characters" :value="knownCharacters">
              <template #suffix>
                %
              </template>
            </n-statistic>
            <n-statistic label="Total Words" :value="totalWords" />

            <n-table :bordered="false" :single-line="false">
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
                  :value="pair.words">
                  <td>{{pair.target}}</td>
                  <td>{{pair.words}}</td>
                </tr>
              </tbody>
            </n-table>
          </n-tab-pane>
          <n-tab-pane name="UnknownWords" tab="View Unknown Words">
            <unknown-words
              class="h-96"
              :words="words"
              :bookFilter="bookId"
            />
          </n-tab-pane>
        </n-tabs>
      </n-layout-content>
      <n-layout-footer class="p-4" bordered>
        <n-space justify="end">

          <n-button type="primary" @click="makeFlashCards">
            Make flash cards
          </n-button>
        </n-space>
      </n-layout-footer>
    </n-layout>
  </n-layout>
</template>

<script lang="ts" setup>
import UnknownWords from '@components/UnknownWords.vue';
import { provide } from 'vue';
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent,
  NLayoutFooter, NButton, NSpace, NTabs, NTabPane,
  NStatistic, NTable,
} from 'naive-ui';
import { invoke } from '@tauri-apps/api/tauri';
import { convertFileSrc } from '@tauri-apps/api/tauri';
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';
import type { Book } from '@/shared/types';

const props = defineProps({
  bookId: {
    required: true,
    type: Number,
  },
});

provide('preferBook', props.bookId);

const book:Book = await window.nodeIpc.loadBook(props.bookId);
const assetUrl = await convertFileSrc(book.cover);
const words = await invoke('learning_target', {
  bookIds: [book.bookId],
});
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
  const topWords: string[] = (
    await window.nodeIpc.topUnknownWords(props.bookId, 50)
  );
  topWords.forEach((word) => {
    store.addWord(word, ActionsEnum.CREATE, { preferBook: props.bookId });
  });
}
</script>
