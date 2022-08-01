<template>
  <div class="container mx-auto px-4">
    <div class="text-center">
      <h2 class="text-center mt-5">About</h2>
      <p>Designed to help manage a dynamic flashcard library to aid in reading books</p>
    </div>
  </div>
</template>

<script>

export default {
  name: 'AboutPage',
  components: {
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
