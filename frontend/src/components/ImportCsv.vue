<template>
  <div>
    <button
      class="btn-primary btn w-full"
      @click="() => importCsvModal = true"
    >
      Import from CSV
    </button>

    <Teleport to="body">
      <div
        :class="['modal', {'modal-open': importCsvModal}]"
        @click="() => importCsvModal = false"
      >
        <div
          class="modal-box w-1/2 max-w-5xl"
          @click.stop
        >
          <div>
            Requires a csv file without a header and each line a new word.

            The second column (if exists) will be parsed the size of the
            anki interval for this word (in days)
          </div>
          <button class="btn-secondary btn" @click="importCsv">
            Import
          </button>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script lang="ts" setup>
import {
  ImportCSVWords,
} from '@wailsjs/backend/KnownWords';
import { FilePicker } from '@wailsjs/main/App';
import { ref } from 'vue';

const importCsvModal = ref(false);

function importCsv() {
  FilePicker('csv')
    .then((csvPath) => {
      ImportCSVWords(csvPath).catch((err) => {
        console.log('Failed to parse file', err);
      });
    })
    .catch((err) => {
      console.log('Failed to open file', err);
    });
}
</script>
