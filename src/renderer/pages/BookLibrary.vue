<template>
  <div class="text-center">
    <h2 class="text-center text-xl mt-5">Your Library</h2>
    <p>Click on a book to start making flashcards.</p>
    <settings-checkbox
      :setting="UserSettings.BookLibrary.OnlyFavorites"
      @update="updateFilter"
    />
    <n-button @click="syncCalibre">Sync Calibre</n-button>
  </div>

  <n-grid x-gap="12" y-gap="12" :cols="4" v-if="books.length > 0">
    <n-gi
      v-for="book in favoriteFilter"
      :key="book.bookId">
      <book-card class="h-[700px]" :book="book" />
    </n-gi>
  </n-grid>
</template>

<script lang="ts" setup>
import {
  onBeforeMount, inject, ref, computed, Ref,
} from 'vue';
import { NGrid, NGi, NButton } from 'naive-ui';
import BookCard from '@/components/BookCard.vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import { Book, UserSettingsKey, UserSettingsType } from '../../shared/types';

function syncCalibre() {
  window.ipc.importCalibreBooks();
}

const UserSettings = inject(UserSettingsKey) as UserSettingsType;

// TODO Would be nice if these properties them selves were reactive
const onlyFavorites = ref(UserSettings.BookLibrary.OnlyFavorites.read());
function updateFilter() {
  onlyFavorites.value = UserSettings.BookLibrary.OnlyFavorites.read();
}
const books: Ref<Book[]> = ref([]);

const favoriteFilter = computed(
  () => books.value
    .filter((book:Book) => {
      if (!onlyFavorites.value) return true;
      return book.favorite;
    }).sort((bookA, bookB) => {
      const aKnown = (bookA.stats.totalKnownWords / bookA.stats.totalWords);
      const bKnown = (bookB.stats.totalKnownWords / bookB.stats.totalWords);
      return bKnown - aKnown;
    }),
);
onBeforeMount(async () => {
  books.value = await window.ipc.loadBooks();
});
</script>
