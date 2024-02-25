<template>
  <div>
    <button
      class="
          btn-primary
          btn w-full"
      @click="() => dataManagerModal = true"
    >
      Check Portable
    </button>

    <Teleport to="#app-root">
      <div
        :class="[
          'modal',
          { 'modal-open': dataManagerModal }]"
        @click="() => dataManagerModal = false"
      >
        <div class="modal-box w-1/2 max-w-5xl" @click.stop>
          <div
            v-if="containsAbsPath"
            class="m-4 rounded-xl bg-warning p-4"
          >
            <div>
              We have detected that some of the data has been saved in a
              non-portable way. While I have automated the process of converting
              it to be made portable. This involves altering data in the
              databases and there is always a chance of there being errors in
              this code.

              To be extra safe you can make a copy of the folder
              <span class="bg-primary font-mono"> {{ configDir }} </span> before
              proceeding
            </div>
            <div
              v-if="badBooks.length > 0"
              class="m-4 rounded-xl bg-error p-4"
            >
              There was a problem with the following books
              <ul class="m-4 list-disc pl-4">
                <li v-for="book in badBooks" :key="book.bookId">
                  {{ book.author }} - {{ book.title }}
                </li>
              </ul>
              Please delete them from read-chinese and reimport them to fix
              the problem
            </div>
            <button class="btn-secondary btn" @click="fixPaths">
              Fix paths
            </button>
          </div>
          <div v-else class="flex gap-4">
            <p>
              Your data should be good to go. You should be able to copy the
              folder <span class="bg-primary font-mono"> {{ configDir }} </span>
              to the equivalent location on your new pc

              If you need any help, please reach out in the discord
            </p>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script lang="ts" setup>
import {
  BookPathsPortable,
  FixBookPaths,
} from '@wailsjs/backend/bookLibrary';
import {
  GetConfigDir,
} from '@wailsjs/backend/ConfigExposer';
import {
  SettingsPathsPortable,
  FixSettingsPaths,
} from '@wailsjs/backend/UserSettings';
import { ref } from 'vue';
import type { backend } from '@wailsjs/models';

const dataManagerModal = ref(false);
const containsAbsPath = ref(false);
const configDir = ref('');
const badBooks = ref<backend.Book[]>([]);
async function load() {
  const booksOk = await BookPathsPortable();
  const metaOk = await SettingsPathsPortable();
  configDir.value = await GetConfigDir();
  containsAbsPath.value = !(metaOk && booksOk);
}

load();

async function fixPaths() {
  // Catch here.
  await FixBookPaths().then(
    (theBadBooks: backend.Book[]) => {
      badBooks.value = theBadBooks;
      load();
    });
  await FixSettingsPaths();
  load();
}
</script>
