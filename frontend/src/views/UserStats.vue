<template>
  <div ref="styleRef" class="flex">
    <div class="w-1/5">
      <div class="flex flex-col p-4">
        <div
          class="stats stats-vertical bg-primary
        text-primary-content"
        >
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
    </div>
    <div class="h-full w-4/5 p-4">
      <Line v-if="loaded" :data="data" :options="options" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { GetWordStats, GetStatsInfo } from '@wailsjs/backend/KnownWords';
import { TotalRead, TotalReadChars } from '@wailsjs/backend/bookLibrary';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import type { ChartData } from 'chart.js';
import { Line } from 'vue-chartjs';
// import type { ChartProps } from 'vue-chartjs';
import { ref, onMounted } from 'vue';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
);

const styleRef = ref<HTMLElement | null>(null);

const totalWords = ref(0);
const totalChars = ref(0);
const words = ref(0);
const characters = ref(0);
const data = ref<ChartData<'line'>>({
  datasets: [],
});
const loaded = ref(false);

onMounted(async () => {
  if (!styleRef.value) {
    return;
  }
  totalWords.value = await TotalRead();
  totalChars.value = await TotalReadChars();

  const style = getComputedStyle(styleRef.value);
  const primary = `hsl(${style.getPropertyValue('--p')}`;
  const secondary = `hsl(${style.getPropertyValue('--s')}`;

  const wordStats = await GetWordStats();
  words.value = wordStats.words;
  characters.value = wordStats.characters;

  const rawdata = await GetStatsInfo();
  data.value = {
    labels: rawdata.map(d => d.day),
    datasets: [{
      label: 'Known Words',
      backgroundColor: primary,
      data: rawdata.map(d => d.known),
    },
    {
      label: 'Known Chars',
      backgroundColor: secondary,
      data: rawdata.map(d => d.knownCharacters),
    },
    ],
  };
  loaded.value = true;
});

const options = {
  responsive: true,
};

</script>
