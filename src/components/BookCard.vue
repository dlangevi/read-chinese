<template>
  <n-card :title="book.title" class="p-4">
    <template #cover>
      <img @click="bookBigMode"
        class="rounded rounded-t h-full w-auto"
        :src="'data:image/png;base64,' + book.imgData"
        :alt="book.title"
      />
    </template>
    <p>{{ book.author }}</p>
    <small>
      ID {{book.bookId}}
    </small>
  </n-card>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { NCard } from 'naive-ui';

const router = useRouter();
const props = defineProps({
  book: {
    type: Object,
    required: true,
  },
});

function bookBigMode() {
  router.push(`/book/${props.book.bookId}`);
}
</script>
