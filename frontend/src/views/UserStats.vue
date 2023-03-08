<template>
  <div class="flex">
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
      <Line class="" :data="data" :options="options" />
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
import { Line } from 'vue-chartjs';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
);

const totalWords = await TotalRead();
const totalChars = await TotalReadChars();
const {
  words,
  characters,
} = await GetWordStats();

const rawdata = await GetStatsInfo();
const data = {
  labels: rawdata.map(d => d.day),
  datasets: [{
    label: 'Known Words',
    backgroundColor: '#f87979',
    data: rawdata.map(d => d.known),
  },
  {
    label: 'Known Chars',
    backgroundColor: '#002200',
    data: rawdata.map(d => d.knownCharacters),
  },
  ],
};
const options = {
  responsive: true,
};

</script>
