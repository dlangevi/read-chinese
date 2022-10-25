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

const wordStats = await window.nodeIpc.wordStats();
const readWords = await window.nodeIpc.totalRead();

const { totalWords } = readWords;
const {
  words,
  characters,
} = wordStats;

</script>
