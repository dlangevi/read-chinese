<template>
  <div
    :class="[
      'card shadow-xl p-4 flex',
      {'bg-success text-success-content': (known > 90)},
      { 'bg-base-200': (known <= 90) },
    ]"
  >
    <figure>
      <img
        v-if="book.cover"
        class="m-auto max-h-full grow rounded object-contain"
        :src="'/' + book.cover"
        :alt="book.title"
        @click="bookBigMode"
      >
      <img
        v-else
        class="m-auto max-h-full grow rounded object-contain"
        src="../assets/empty-book.jpg"
        alt="Default image"
        @click="bookBigMode"
      >
    </figure>
    <div class="card-body">
      <h2 class="card-title">
        {{ `${book.title} - ${book.author}` }}
      </h2>
      <div class="m-4 text-center">
        Known: {{ known.toFixed(2) }}%
      </div>
      <div class="flex flex-wrap place-content-end gap-2">
        <button
          v-if="!isFavorite"
          class="btn-accent btn-xs btn"
          @click="favorite(true)"
        >
          Add To Favorites
        </button>
        <button
          v-if="isFavorite"
          class="btn-accent btn-xs btn"
          @click="favorite(false)"
        >
          Remove From Favorites
        </button>
        <button
          v-if="!isRead"
          class="btn-accent btn-xs btn"
          @click="markRead(true)"
        >
          Mark Read
        </button>
        <button
          v-if="isRead"
          class="btn-accent btn-xs btn"
          @click="markRead(false)"
        >
          Unmark Read
        </button>
        <button
          class="btn-accent btn-xs btn"
          @click="deleteBook"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router';
import { ref } from 'vue';

import { DeleteBook, SetFavorite, SetRead } from '@wailsjs/backend/bookLibrary';

const router = useRouter();
const props = defineProps({
  book: {
    type: Object,
    required: true,
  },
});

const isRead = ref(props.book.hasRead);
const isFavorite = ref(props.book.favorite);
function favorite(setTo:boolean) {
  SetFavorite(props.book.bookId, setTo);
  isFavorite.value = setTo;
}
function markRead(setTo:boolean) {
  SetRead(props.book.bookId, setTo);
  isRead.value = setTo;
}

function deleteBook() {
  DeleteBook(props.book.bookId);
}

const known = (
  (props.book.stats.totalKnownWords / props.book.totalWords) * 100
);

function bookBigMode() {
  router.push(`/book/${props.book.bookId}`);
}
</script>
