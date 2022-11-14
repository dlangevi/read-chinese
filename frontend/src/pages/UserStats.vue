<template>
  <div>
    Look at these stats!
    <div
      class="
        grid
        grid-cols-2
        gap-4
        place-content-center
        h-full
        w-1/2
        p-16
        ">
      <n-statistic label="Known Words" :value="words" />
      <n-statistic label="Known Characters" :value="characters" />
      <n-statistic label="Total words read" :value="totalWords" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { NStatistic } from 'naive-ui';
import { GetWordStats } from '@wailsjs/backend/KnownWords';
import { TotalRead } from '@wailsjs/backend/bookLibrary';

const totalWords = await TotalRead();
const wordStats = await GetWordStats();
const {
  words,
  characters,
} = wordStats;

</script>
