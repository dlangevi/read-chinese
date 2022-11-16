<template>
  <n-card
    :title="`${book.title} - ${book.author}`"
    :class="(known > 90) ? 'bg-green-300 p-4' : 'p-4'"
  >
    <template #cover>
      <img
        class="m-auto max-h-full w-auto rounded"
        :src="book.cover"
        :alt="book.title"
        @click="bookBigMode"
      >
    </template>
    <div class="m-4 text-center">
      Known: {{ known.toFixed(2) }}%
    </div>
    <n-space justify="end">
      <button
        v-if="!isFavorite"
        class="btn btn-xs btn-accent"
        @click="favorite(true)"
      >
        Add To Favorites
      </button>
      <button
        v-if="isFavorite"
        class="btn btn-xs btn-accent"
        @click="favorite(false)"
      >
        Remove From Favorites
      </button>
      <button
        v-if="!isRead"
        class="btn btn-xs btn-accent"
        @click="markRead(true)"
      >
        Mark Read
      </button>
      <button
        v-if="isRead"
        class="btn btn-xs btn-accent"
        @click="markRead(false)"
      >
        Unmark Read
      </button>
      <button
        class="btn btn-xs btn-accent"
        @click="deleteBook"
      >
        Delete
      </button>
    </n-space>
  </n-card>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router';
import { NSpace, NCard } from 'naive-ui';
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
  (props.book.stats.totalKnownWords / props.book.stats.totalWords) * 100
);

function bookBigMode() {
  router.push(`/book/${props.book.bookId}`);
}
</script>
