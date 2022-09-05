<template>
  <n-layout class="m-4" has-sider>
    <n-layout-sider width=500 bordered content-style="padding: 24px;">
      <p>{{ book.author }}</p>
      <small>
        ID {{book.bookID}}
      <img
        class="rounded rounded-t w-auto"
        :src="'data:image/png;base64,' + book.imgData"
        :alt="book.title"
      />
      </small>
    </n-layout-sider>
    <n-layout>
      <n-layout-header class="p-4" bordered>
        Stats
      </n-layout-header>
      <n-layout-content class="p-8">
        Pingshan Road
      </n-layout-content>
      <n-layout-footer class="p-4" bordered>
        Chengfu Road
      </n-layout-footer>
    </n-layout>
  </n-layout>
</template>

<script setup>
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent,
  NLayoutFooter,
} from 'naive-ui';

const props = defineProps({
  bookID: String,
});

const book = await window.ipc.loadBook(props.bookID);
</script>
