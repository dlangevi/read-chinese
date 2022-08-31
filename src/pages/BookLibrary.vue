<template>
  <div class="container mx-auto px-4">
    <div class="text-center">
      <h2 class="text-center text-xl mt-5">Your Library</h2>
      <p>Click on a book to start making flashcards.</p>
    </div>

    <div class="grid grid-cols-3" v-if="books.length > 0">
        <book-card :book="book" v-for="(book, i) in books" :key="i" />
    </div>
  </div>
</template>

<script>
import BookCard from '../components/BookCard.vue';

export default {
  name: 'BookLibrary',
  components: {
    BookCard,
  },
  data() {
    return {
      books: [],
    };
  },
  methods: {
  },
  async beforeRouteEnter(to, from, next) {
    const books = await window.ipc.loadBooks();
    next((vm) => { vm.books = books; });
  },
};
</script>
