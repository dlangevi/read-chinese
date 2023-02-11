<template>
  <with-sidebar>
    <template #sidebar>
      <settings-checkbox
        :setting="UserSettings.BookLibrary.OnlyFavorites"
        @update="updateFilter"
      />
      <settings-checkbox
        :setting="UserSettings.BookLibrary.HideRead"
        @update="updateFilter"
      />
      <button
        class="btn-primary btn-sm btn"
        @click="syncCalibre"
      >
        Sync Calibre
      </button>
      <button
        class="btn-primary btn-sm btn"
        @click="exportBooks"
      >
        Export Book Stats
      </button>
    </template>
    <div class="text-center">
      <h2 class="mt-5 text-xl">
        Your Library
      </h2>
      <p>Click on a book to start making flashcards.</p>
    </div>
    <div
      v-if="books.length > 0"
      class="m-8 grid grid-cols-4 gap-12"
    >
      <div
        v-for="book in favoriteFilter"
        :key="book.bookId"
      >
        <book-card
          class="h-[700px]"
          :book="book"
        />
      </div>
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import {
  onBeforeMount, ref, computed, Ref,
} from 'vue';
import BookCard from '@/components/BookCard.vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import { getUserSettings } from '@/lib/userSettings';
import { backend } from '@wailsjs/models';
import { SaveFile } from '@wailsjs/main/App';
import { GetBooks, GetDetailedBooks } from '@wailsjs/backend/bookLibrary';
import { ImportCalibreBooks } from '@wailsjs/backend/Calibre';
import WithSidebar from '@/layouts/WithSidebar.vue';

async function syncCalibre() {
  const err = await ImportCalibreBooks();
  console.log(err);
}

async function exportBooks() {
  const filename = await SaveFile();
  const detailed = await GetDetailedBooks(filename);
  console.log(JSON.stringify(detailed));
}

const UserSettings = getUserSettings();

// TODO Would be nice if these properties them selves were reactive
const onlyFavorites = ref(UserSettings.BookLibrary.OnlyFavorites.read());
const hideRead = ref(UserSettings.BookLibrary.HideRead.read());
function updateFilter() {
  onlyFavorites.value = UserSettings.BookLibrary.OnlyFavorites.read();
  hideRead.value = UserSettings.BookLibrary.HideRead.read();
}
const books: Ref<backend.Book[]> = ref([]);

const favoriteFilter = computed(
  () => books.value
    .filter((book:backend.Book) => {
      if (hideRead.value && book.hasRead) {
        return false;
      }
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
