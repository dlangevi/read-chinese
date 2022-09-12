<template>
  <n-layout class="m-4" has-sider>
    <n-layout-sider width=500 bordered content-style="padding: 24px;">
      <img
        class="rounded rounded-t w-auto"
        :src="'data:image/png;base64,' + book.imgData"
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
            Known: {{known}}%
            Can Read: {{likelyKnown}}
            Characters Known: {{knownCharacters}}
          </n-tab-pane>
          <n-tab-pane name="UnknownWords" tab="View Unknown Words">
            <unknown-words
              class="h-96"
              :bookFilter="[parseInt(props.bookId)]"
              />
          </n-tab-pane>
        </n-tabs>

      </n-layout-content>
      <n-layout-footer class="p-4" bordered>
        <n-space justify="end">
          <n-button type="primary">
            Mark words known
          </n-button>
          <n-button type="primary">
            Make flash cards
          </n-button>
        </n-space>
      </n-layout-footer>
    </n-layout>
  </n-layout>
</template>

<script setup>
import UnknownWords from '@/components/UnknownWords.vue';
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent,
  NLayoutFooter, NButton, NSpace, NTabs, NTabPane,
} from 'naive-ui';

const props = defineProps({
  bookId: String,
});

const book = await window.ipc.loadBook(props.bookId);

const known = ((book.totalKnownWords / book.totalWords) * 100).toFixed(2);
const likelyKnown = (
  (book.probablyKnownWords / book.totalWords) * 100).toFixed(2);
const knownCharacters = (
  (book.knownCharacters / book.totalCharacters) * 100).toFixed(2);
console.log(book);
</script>
