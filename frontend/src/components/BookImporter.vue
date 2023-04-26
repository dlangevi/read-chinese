<template>
  <div>
    <button
      class="btn-primary btn w-full"
      @click="openImporter"
    >
      Sync Books From Calibre
    </button>
    <Teleport to="#app-root">
      <div
        :class="['modal', {'modal-open': importBooksModal}]"
        @click="() => importBooksModal = false"
      >
        <div
          class="modal-box relative w-1/2 max-w-5xl"
          @click.stop
        >
          <div class="overflow-x-auto">
            <div v-if="newBooks.length === 0">
              All books imported
            </div>
            <div v-else class="flex flex-col gap-4">
              <h2 class="mx-4 text-3xl">Unimported books</h2>
              <table class="table w-full">
                <thead>
                  <tr>
                    <th />
                    <th>Title</th>
                    <th>Author</th>
                    <th>Formats</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="book in newBooks"
                    :key="book.id"
                    class="hover"
                  >
                    <td :id="'bookId' + book.id">
                      <label>
                        <input
                          v-model="selectedBooks"
                          type="checkbox"
                          class="checkbox"
                          :value="book"
                        >
                      </label>
                    </td>
                    <td>{{ book.title }}</td>
                    <td>{{ book.authors }}</td>
                    <td>{{ getFormats(book.formats) }}</td>
                  </tr>
                </tbody>
              </table>
              <div class="flex gap-4">
                <button
                  class="btn-primary btn grow basis-1/2"
                  @click="syncCalibre(selectedBooks)"
                >
                  Import Selected Books
                </button>
                <button
                  class="btn-primary btn grow basis-1/2"
                  @click="syncCalibre(newBooks)"
                >
                  Import All Books
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useLoader } from '@/lib/loading';
import { ImportCalibreBooks, GetCalibreBooks } from '@wailsjs/backend/Calibre';
import { backend } from '@wailsjs/models';

function getFormats(formats : string[]) {
  return formats.map((path) => path.split('.').slice(-1)).join(',');
}

const loader = useLoader();

const books = ref<backend.CalibreBook[]>([]);
const selectedBooks = ref<backend.CalibreBook[]>([]);
const importBooksModal = ref(false);

const newBooks = computed(() => {
  return books.value.filter(book => !book.exists);
});

async function syncCalibre(books : backend.CalibreBook[]) {
  return loader.withLoader(async () => {
    await ImportCalibreBooks(books);
    importBooksModal.value = false;
  });
}
async function getCalibreBooks() {
  books.value = await GetCalibreBooks();
}
getCalibreBooks();

async function openImporter() {
  await getCalibreBooks();
  importBooksModal.value = true;
}

</script>
