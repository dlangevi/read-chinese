<template>
  <div class="text-center">
    <h2 class="text-center text-xl mt-5">
      Your Library
    </h2>
    <p>Click on a book to start making flashcards.</p>
    <settings-checkbox
      :setting="UserSettings.BookLibrary.OnlyFavorites"
      @update="updateFilter"
    />
    <n-button @click="syncCalibre">
      Sync Calibre
    </n-button>
  </div>

  <n-grid
    v-if="books.length > 0"
    x-gap="12"
    y-gap="12"
    :cols="4"
  >
    <n-gi
      v-for="book in favoriteFilter"
      :key="book.bookId"
    >
      <book-card
        class="h-[700px]"
        :book="book"
      />
    </n-gi>
  </n-grid>
</template>

<script lang="ts" setup>
import {
  onBeforeMount, ref, computed, Ref,
} from 'vue';
import { NGrid, NGi, NButton } from 'naive-ui';
import BookCard from '@/components/BookCard.vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import { getUserSettings } from '@/lib/userSettings';
import { backend } from '@wailsjs/models';
import { GetBooks } from '@wailsjs/backend/bookLibrary';
import { ImportCalibreBooks } from '@wailsjs/backend/Calibre';

async function syncCalibre() {
  const err = await ImportCalibreBooks();
  console.log(err);
}

const UserSettings = getUserSettings();

// TODO Would be nice if these properties them selves were reactive
const onlyFavorites = ref(UserSettings.BookLibrary.OnlyFavorites.read());
function updateFilter() {
  onlyFavorites.value = UserSettings.BookLibrary.OnlyFavorites.read();
}
const books: Ref<backend.Book[]> = ref([]);

const favoriteFilter = computed(
  () => books.value
    .filter((book:backend.Book) => {
      if (!onlyFavorites.value) return true;
      return book.favorite;
    }).sort((bookA, bookB) => {
      const aKnown = (bookA.stats.totalKnownWords / bookA.stats.totalWords);
      const bKnown = (bookB.stats.totalKnownWords / bookB.stats.totalWords);
      return bKnown - aKnown;
    }),
);
onBeforeMount(async () => {
  books.value = await GetBooks();
  console.log(books.value);
});
</script>
