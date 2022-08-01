<template>
  <div class="container mx-auto px-4">
    <div class="text-center">
      <h2 class="text-center mt-5">Your Library</h2>
      <p>You clicked on {{ bookID }}.</p>
      <book-card :book="book"/>
    </div>
  </div>
</template>

<script>
import BookCard from '../components/BookCardTail.vue';

export default {
  name: 'BookStats',
  components: {
    BookCard,
  },
  props: {
    bookID: String,
  },
  data() {
    return {
      book: [],
    };
  },
  methods: {
    async getBook() {
      this.book = await window.ipc.loadBook(this.bookID);
      console.log(this.book);
    },
  },
  mounted() {
    this.getBook();
  },
};
</script>
