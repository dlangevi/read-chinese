<template>
  <main>
    <div class="card mb-4 shadow-lg border border-secondary bg-dark text-light">
      <img
        class="card-img-top"
        :src="'data:image/png;base64,' + book.imgData"
        alt="Card image cap"
      />
      <div class="card-body">
        <p class="card-text">{{ book.title }}</p>
        <p class="card-text">{{ book.author }}</p>
        <div class="">
          <div class="d-flex justify-content-between align-items-center">
            <div class="btn-group">
              <button type="button" class="btn btn-sm btn-outline-secondary">
                <i class="fas fa-heart"></i>
              </button>
              <button type="button" class="btn btn-sm btn-outline-secondary">
                <i class="fas fa-plus"></i>
              </button>
            </div>
            <small class="ms-auto text-muted"
              >ID {{book.id}}</small
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
};
</script>
