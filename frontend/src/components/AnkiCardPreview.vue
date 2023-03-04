<template>
  <div class="card flex flex-col gap-4 border-4 border-base-300 p-4">
    <h2 class="text-xl font-bold">
      Hanzi: {{ cardManager.word }}
    </h2>
    <div>
      <h2 class="text-l inline font-bold">
        Sentence:
      </h2>
      <span> {{ cardManager.sentence }}</span>
    </div>
    <div @click="cardManager.changeStep(StepsEnum.ENGLISH)">
      <h2 class="text-l inline font-bold">
        Definition:
      </h2>
      <!-- Definitions can have html formatting in them -->
      <p>
        <span
          v-for="definition in cardManager.englishDefn"
          :key="definition.definition"
        >
          {{ formatDefinition(definition) }}
          <br>
        </span>
      </p>
    </div>
    <div v-if="cardManager.steps.includes(StepsEnum.CHINESE)">
      <h2 class="text-l inline font-bold">
        ChineseDefinition:
      </h2>
      <p>
        <span
          v-for="definition in cardManager.chineseDefn"
          :key="definition.definition"
        >
          {{ formatDefinition(definition) }}
        </span>
      </p>
    </div>
    <div v-if="cardManager.steps.includes(StepsEnum.IMAGE)">
      <h2 class="text-l inline font-bold">
        Images:
      </h2>
      <div class="flex flex-wrap gap-1">
        <img
          v-for="image, i in cardManager.images"
          :key="i"
          class="max-h-24 w-auto"
          :src="getImageSrc(image)"
          :alt="image.name ||
            'Image related to search word, no alt text generated'"
        >
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import {
  useCardManager,
  formatDefinition,
  getImageSrc,
} from '@/stores/CardManager';

const cardManager = useCardManager();

defineEmits(['change-step']);

</script>
