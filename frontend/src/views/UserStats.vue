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
import { GetWordStats } from '@wailsjs/backend/KnownWords';
import { TotalRead, TotalReadChars } from '@wailsjs/backend/bookLibrary';

const totalWords = await TotalRead();
const totalChars = await TotalReadChars();
const wordStats = await GetWordStats();
const {
  words,
  characters,
} = wordStats;

</script>
