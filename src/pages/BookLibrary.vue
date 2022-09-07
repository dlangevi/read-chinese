<template>
  <div class="text-center">
    <h2 class="text-center text-xl mt-5">Your Library</h2>
    <p>Click on a book to start making flashcards.</p>
  </div>

  <n-grid x-gap="12" y-gap="12" :cols="4" v-if="books.length > 0">
    <n-gi v-for="(book, i) in books" :key="i">
      <book-card :book="book" />
    </n-gi>
  </n-grid>
</template>

<script setup>
import { onBeforeMount, ref } from 'vue';
import { NGrid, NGi } from 'naive-ui';
import BookCard from '@/components/BookCard.vue';

const books = ref([]);
onBeforeMount(async () => {
  books.value = await window.ipc.loadBooks();
});
</script>
