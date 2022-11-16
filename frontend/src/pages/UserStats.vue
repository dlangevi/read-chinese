<template>
  <div class="m-auto flex h-full w-1/4">
    <div class="stats my-auto bg-primary text-primary-content shadow">
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { GetWordStats } from '@wailsjs/backend/KnownWords';
import { TotalRead } from '@wailsjs/backend/bookLibrary';

const totalWords = await TotalRead();
const wordStats = await GetWordStats();
const {
  words,
  characters,
} = wordStats;

</script>
