<template>
  <with-sidebar>
    <template #sidebar>
      <settings-checkbox
        :setting="ComponentTable.BookLibrary.OnlyFavorites"
        :initial-value="UserSettings.BookLibrary.OnlyFavorites"
      />

      <settings-checkbox
        :setting="ComponentTable.BookLibrary.HideRead"
        :initial-value="UserSettings.BookLibrary.HideRead"
      />
      <settings-checkbox
        :setting="ComponentTable.BookLibrary.DisplayTable"
        :initial-value="UserSettings.BookLibrary.DisplayTable"
      />
      <book-importer />
      <button
        class="btn-primary btn"
        @click="exportBooks"
      >
        Export Book Stats
      </button>
      <button
        class="btn-primary btn"
        @click="recalculateBooks"
      >
        Resegment Books
      </button>
    </template>
    <div class="text-center">
      <h2 class="mt-5 text-xl">
        Your Library
      </h2>
      <p>Click on a book to start making flashcards.</p>
    </div>
    <div v-if="books.length > 0">
      <book-table
        v-if="UserSettings.BookLibrary.DisplayTable"
        :books="favoriteFilter"
      />

      <div
        v-else
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
    </div>
  </with-sidebar>
</template>

<script lang="ts" setup>
import {
  onUnmounted, onBeforeMount, ref, computed, Ref,
} from 'vue';
import BookCard from '@/components/BookCard.vue';
import BookTable from '@/components/BookTable.vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import BookImporter from '@/components/BookImporter.vue';
import { getUserSettings, ComponentTable } from '@/lib/userSettings';
import { backend } from '@wailsjs/models';
import { SaveFile } from '@wailsjs/backend/Backend';
import {
  GetDetailedBooks,
  RecalculateBooks, ExportDetailedBooks,
} from '@wailsjs/backend/bookLibrary';
import { useLoader } from '@/lib/loading';
import WithSidebar from '@/layouts/WithSidebar.vue';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
const loader = useLoader();

async function exportBooks() {
  const filename = await SaveFile();
  return ExportDetailedBooks(filename);
}

async function recalculateBooks() {
  return loader.withLoader(RecalculateBooks);
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
        const aKnown = (bookA.stats.totalKnownWords / bookA.totalWords);
        const bKnown = (bookB.stats.totalKnownWords / bookB.totalWords);
        return bKnown - aKnown;
      });
  });

onBeforeMount(async () => {
  books.value = await GetDetailedBooks();
  console.log(books.value);
  EventsOn('BooksUpdated', (newBooks : backend.Book[]) => {
    books.value = newBooks;
  });
});

onUnmounted(async () => {
  EventsOff('BooksUpdated');
});
</script>
