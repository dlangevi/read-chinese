<template>
  <main>
    <div class="rounded border m-4
      bg-white border-3 mb-4 border border-grey" @click="bookBigMode">
      <img
        class="w-full rounded rounded-t"
        :src="'data:image/png;base64,' + book.imgData"
        alt="Card image cap"
      />
      <div class="p-6">
        <p class="mb-0">{{ book.title }}</p>
        <p class="card-text">{{ book.author }}</p>
        <div class="">
          <div class="flex justify-content-between align-items-center">
            <div class="relative inline-flex align-middle">
              <button type="button" class="inline-block align-middle
                text-center select-none border font-normal whitespace-no-wrap
                py-2 px-4 rounded text-base leading-normal no-underline py-1 px-2
                text-sm leading-tight text-grey-dark border-grey bg-white
                hover:bg-grey-light hover:text-grey-darker">
                <i class="fas fa-heart"></i>
              </button>
              <button type="button" class="btn btn-sm btn-outline-secondary">
                <i class="fas fa-plus"></i>
              </button>
            </div>
            <small class="ms-auto text-grey"
              >ID {{book.bookID}}</small
            >
          </div>
        </div>
      </div>
    </div>
  </main>
</template>
<script>
export default {
  name: 'BookCard',
  props: {
    book: {
      type: Object,
      required: true,
    },
  },
  methods: {
    bookBigMode() {
      this.$router.push(`/book/${this.book.bookID}`);
    },
  },
};
</script>
