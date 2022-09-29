<template>
  <div class="text-center">
    <h2 class="text-center text-xl mt-5">Your Library</h2>
    <p>Click on a book to start making flashcards.</p>
    <settings-checkbox
      :setting="UserSettings.BookLibrary.OnlyFavorites"
      @update="updateFilter"
    />
  </div>

  <n-grid x-gap="12" y-gap="12" :cols="4" v-if="books.length > 0">
    <n-gi v-for="(book, i) in favoriteFilter" :key="i">
      <book-card class="h-[700px]" :book="book" />
    </n-gi>
  </n-grid>
</template>

<script setup>
import {
  onBeforeMount, inject, ref, computed,
} from 'vue';
import { NGrid, NGi } from 'naive-ui';
import BookCard from '@/components/BookCard.vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';

const UserSettings = inject('userSettings');

// TODO Would be nice if these properties them selves were reactive
const onlyFavorites = ref(UserSettings.BookLibrary.OnlyFavorites.read());
function updateFilter() {
  onlyFavorites.value = UserSettings.BookLibrary.OnlyFavorites.read();
}
const books = ref([]);

const favoriteFilter = computed(() => books.value.filter((book) => {
  if (!onlyFavorites.value) return true;
  return book.favorite;
}));
onBeforeMount(async () => {
  books.value = await window.ipc.loadBooks();
});
</script>
