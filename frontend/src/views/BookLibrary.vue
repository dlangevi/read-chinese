<template>
  <with-sidebar>
    <template #sidebar>
      <settings-checkbox
        :setting="ComponentTable.OnlyFavorites"
        :initial-value="UserSettings.BookLibrary.OnlyFavorites"
      />

      <settings-checkbox
        :setting="ComponentTable.HideRead"
        :initial-value="UserSettings.BookLibrary.HideRead"
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
import { getUserSettings, ComponentTable } from '@/lib/userSettings';
import { backend } from '@wailsjs/models';
import { SaveFile } from '@wailsjs/main/App';
import { GetBooks, GetDetailedBooks } from '@wailsjs/backend/bookLibrary';
import { ImportCalibreBooks } from '@wailsjs/backend/Calibre';
import { useLoader } from '@/lib/loading';
import WithSidebar from '@/layouts/WithSidebar.vue';
const loader = useLoader();

async function syncCalibre() {
  return loader.withLoader(ImportCalibreBooks, 'Importing calibre');
}

async function exportBooks() {
  const filename = await SaveFile();
  const detailed = await GetDetailedBooks(filename);
  console.log(JSON.stringify(detailed));
}

const UserSettings = getUserSettings();

const books: Ref<backend.Book[]> = ref([]);

const favoriteFilter = computed(
  () => {
    const onlyFavorites = UserSettings.BookLibrary.OnlyFavorites;
    const hideRead = UserSettings.BookLibrary.HideRead;
    return books.value
      .filter((book:backend.Book) => {
        if (hideRead && book.hasRead) {
          return false;
        }
        if (!onlyFavorites) return true;
        return book.favorite;
      }).sort((bookA, bookB) => {
        const aKnown = (bookA.stats.totalKnownWords / bookA.stats.totalWords);
        const bKnown = (bookB.stats.totalKnownWords / bookB.stats.totalWords);
        return bKnown - aKnown;
      });
  });

onBeforeMount(async () => {
  books.value = await GetBooks();
  console.log('books:', books.value);
});
</script>
