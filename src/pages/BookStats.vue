<template>
  <div class="container mx-auto px-4 h-full">
    <book-card :book="book" class="h-full w-1/3"/>
  </div>
</template>

<script setup>
import BookCard from '../components/BookCard.vue';

const props = defineProps({
  bookID: String,
});

const book = await window.ipc.loadBook(props.bookID);
</script>
