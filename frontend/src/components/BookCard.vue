<template>
  <n-card
    :title="`${book.title} - ${book.author}`"
    :class="(known > 90) ? 'bg-green-300 p-4' : 'p-4'">
    <template #cover>
      <img
        class="rounded rounded-t max-h-full w-auto m-auto"
        :src="book.cover"
        :alt="book.title"
        @click="bookBigMode"
      />
    </template>
    <div class="m-4 text-center">
      Known: {{known.toFixed(2)}}%
    </div>
    <n-space justify="end">
      <n-button v-if="!isFavorite" @click="favorite(true)">
        Add To Favorites
      </n-button>
      <n-button v-if="isFavorite" @click="favorite(false)">
        Remove From Favorites
      </n-button>
      <n-button v-if="!isRead" @click="markRead(true)">Mark Read</n-button>
      <n-button v-if="isRead" @click="markRead(false)">Unmark Read</n-button>
      <n-button @click="deleteBook">Delete</n-button>
    </n-space>
  </n-card>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router';
import { NSpace, NCard, NButton } from 'naive-ui';
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
