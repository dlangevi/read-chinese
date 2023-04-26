<template>
  <div>
    <button
      class="btn-primary btn w-full"
      @click="openAdder"
    >
      Add Book
    </button>
    <Teleport to="#app-root">
      <div
        :class="['modal', {'modal-open': addBookModal}]"
        @click="() => addBookModal = false"
      >
        <div
          class="modal-box relative w-1/2 max-w-5xl"
          @click.stop
        >
          <div class="flex flex-col gap-4 overflow-x-auto">
            <h2 class="text-xl">
              Import a book
            </h2>
            <input
              v-model="bookTitle"
              class="input-bordered input"
              type="text"
              placeholder="Book Title"
            >
            <input
              v-model="bookAuthor"
              class="input-bordered input"
              type="text"
              placeholder="Book Author"
            >
            <div>
              <button class="btn-secondary btn" @click="pickFile">
                Select txt file
              </button>
              {{ bookFile }}
            </div>
            <div v-if="ready">
              <button class="btn-secondary btn" @click="importBook">
                Import Book
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
// import { backend } from '@wailsjs/models';
import { FilePicker } from '@wailsjs/backend/Backend';
import { AddBook } from '@wailsjs/backend/bookLibrary';
const addBookModal = ref(false);

const bookTitle = ref('');
const bookAuthor = ref('');
const bookFile = ref('');
const ready = computed(() => {
  return bookAuthor.value.length > 0 &&
         bookTitle.value.length > 0 &&
         bookFile.value.length > 0;
});

async function openAdder() {
  bookTitle.value = '';
  bookAuthor.value = '';
  addBookModal.value = true;
}

async function pickFile() {
  bookFile.value = await FilePicker('txt');
}

async function importBook() {
  // With loader
  await AddBook(bookAuthor.value, bookTitle.value, '', bookFile.value);
  addBookModal.value = false;
}

</script>
