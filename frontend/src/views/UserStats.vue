<template>
  <div class="flex h-full">
    <div class="stats m-auto bg-primary text-primary-content shadow">
      <div class="stat place-items-center">
        <div class="stat-title"> Known Words </div>
        <div class="stat-value"> {{ words }} </div>
      </div>
      <div class="stat place-items-center">
        <div class="stat-title"> Known Characters </div>
        <div class="stat-value"> {{ characters }} </div>
      </div>
      <div class="stat place-items-center">
        <div class="stat-title"> Total Words Read </div>
        <div class="stat-value"> {{ totalWords }} </div>
      </div>
      <div class="stat place-items-center">
        <div class="stat-title"> Total Characters Read </div>
        <div class="stat-value"> {{ totalChars }} </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { GetWordStats } from '@wailsjs/backend/KnownWords';
import { TotalRead, TotalReadChars } from '@wailsjs/backend/bookLibrary';

const totalWords = ref(0);
const totalChars = ref(0);
const words = ref(0);
const characters = ref(0);
async function run() {
  totalWords.value = await TotalRead();
  totalChars.value = await TotalReadChars();
  const wordStats = await GetWordStats();
  words.value = wordStats.words;
  characters.value = wordStats.characters;
}
run();

</script>
