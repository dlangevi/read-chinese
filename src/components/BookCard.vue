<template>
  <n-card :title="book.title"
    :class="(known > 90) ? 'bg-green-300 p-4' : 'p-4'"
    @click="bookBigMode">
    <template #cover>
      <img
        class="rounded rounded-t h-full w-auto"
        :src="'data:image/png;base64,' + book.imgData"
        :alt="book.title"
      />
    </template>
    <p>{{ book.author }}</p>
    <small>
      Known: {{known}}
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

const known = (
  (props.book.totalKnownWords / props.book.totalWords) * 100).toFixed(2);

function bookBigMode() {
  router.push(`/book/${props.book.bookId}`);
}
</script>
