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
          Hello
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
import { ref } from 'vue';

const importCsvModal = ref(false);

function importCsv() {
  ImportCSVWords('');
}
</script>
