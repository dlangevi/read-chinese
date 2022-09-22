<template>
  <n-card :title="`${book.title} - ${book.author}`"
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
      <n-button @click="favorite">Add To Favorites</n-button>
      <n-button @click="markRead">Mark Read</n-button>
      <n-button @click="deleteBook">Delete</n-button>
    </n-space>
  </n-card>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { NSpace, NCard, NButton } from 'naive-ui';

const router = useRouter();
const props = defineProps({
  book: {
    type: Object,
    required: true,
  },
});

function favorite() {

}
function markRead() {

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
