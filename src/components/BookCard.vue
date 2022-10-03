<template>
  <n-card
    :title="`${book.title} - ${book.author}`"
    :class="(known > 90) ? 'bg-green-300 p-4' : 'p-4'">
    <template #cover>
      <img
        class="rounded rounded-t max-h-full w-auto m-auto"
        :srsc="'data:image/png;base64,' + book.imgData"
        :src="'atom:///' + book.cover"
        :alt="book.title"
        @click="bookBigMode"
      />
    </template>
    <div class="m-4 text-center">
      Known: {{known}}%
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

<script setup>
import { useRouter } from 'vue-router';
import { NSpace, NCard, NButton } from 'naive-ui';
import { ref } from 'vue';

const router = useRouter();
const props = defineProps({
  book: {
    type: Object,
    required: true,
  },
});

const isRead = ref(props.book.hasRead);
const isFavorite = ref(props.book.favorite);
function favorite(setTo) {
  window.ipc.setFavorite(props.book.bookId, setTo);
  isFavorite.value = setTo;
}
function markRead(setTo) {
  window.ipc.setRead(props.book.bookId, setTo);
  isRead.value = setTo;
}

function deleteBook() {
  window.ipc.deleteBook(props.book.bookId);
}

const known = (
  (props.book.totalKnownWords / props.book.totalWords) * 100).toFixed(2);

function bookBigMode() {
  router.push(`/book/${props.book.bookId}`);
}
</script>
