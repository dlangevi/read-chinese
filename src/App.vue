<template>
  <div id="app">
    <div class="container">
      <div class="text-center">
        <h2 class="text-center mt-5">Your Library</h2>
        <p>Click on a book to start making flashcards.</p>
      </div>

      <div class="row" v-if="books.length > 0">
        <div class="col-md-3" v-for="(book, i) in books" :key="i">
          <book-card :book="book" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import BookCard from './components/BookCard.vue';
// import { getTimesRan } from './helpers/database';

export default {
  name: 'App',
  components: {
    BookCard,
  },
  data() {
    return {
      books: [],
    };
  },
  methods: {
    getBooks() {
      window.ipc.send('need-books');
    },
  },
  mounted() {
    window.ipc.on('give-books', (books) => {
      console.log(books.map((book) => `${book.author}-${book.title}`));
      this.books = books;
    });
    this.getBooks();
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
