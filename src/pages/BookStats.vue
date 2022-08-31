<template>
  <div class="container mx-auto px-4">
    <div class="text-center">
      <h2 class="text-center mt-5">Your Library</h2>
      <p>You clicked on {{ bookID }}.</p>
      <book-card :book="book"/>
    </div>
  </div>
</template>

<script setup>
import BookCard from '../components/BookCard.vue';

const props = defineProps({
  bookID: String,
});

const book = await window.ipc.loadBook(props.bookID);
</script>
