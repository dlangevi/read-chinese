<template>
  <n-layout class="m-4" has-sider>
    <n-layout-sider width=500 bordered content-style="padding: 24px;">
      <img
        class="rounded rounded-t w-auto"
        :src="'atom:///' + book.cover"
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

<script setup>
import UnknownWords from '@/components/UnknownWords.vue';
import { provide } from 'vue';
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent,
  NLayoutFooter, NButton, NSpace, NTabs, NTabPane,
  NStatistic, NTable,
} from 'naive-ui';
import { useCardQueue, ActionsEnum } from '@/stores/CardQueue';

const props = defineProps({
  bookId: {
    required: true,
    type: String,
  },
});

provide('preferBook', props.bookId);

const book = await window.ipc.loadBook(props.bookId);

const known = ((book.totalKnownWords / book.totalWords) * 100).toFixed(2);
const likelyKnown = (
  (book.probablyKnownWords / book.totalWords) * 100).toFixed(2);
const knownCharacters = (
  (book.knownCharacters / book.totalCharacters) * 100).toFixed(2);
console.log(book);

const { totalWords } = book;
const firstTarget = book.needToKnow.findIndex((n) => n !== 0);
const targets = book.targets.slice(firstTarget, firstTarget + 3);
const needToKnow = book.needToKnow.slice(firstTarget, firstTarget + 3);
const targetPairs = targets.map((e, i) => (
  { target: e, words: needToKnow[i] }));

const store = useCardQueue();
async function makeFlashCards() {
  const words = await window.ipc.topUnknownWords(props.bookId, 50);
  words.forEach(({ word }) => {
    store.addWord(word, ActionsEnum.CREATE, { preferBook: props.bookId });
  });
  console.log(words);
}
</script>
