<template>
  <div class="card">
    <h2 class="text-xl font-bold">
      Hanzi: {{ cardManager.word }}
    </h2>
    <div>
      <button
        class="btn-primary btn-xs btn m-2 inline"
        @click="cardManager.changeStep(StepsEnum.SENTENCE)"
      >
        Edit
      </button>
      <h2 class="text-l inline font-bold">
        Sentence:
      </h2>
      <span> {{ cardManager.sentence }}</span>
    </div>
    <div>
      <button
        class="btn-primary btn-xs btn m-2"
        @click="cardManager.changeStep(StepsEnum.ENGLISH)"
      >
        Edit
      </button>
      <h2 class="text-l inline font-bold">
        Definition:
      </h2>
      <p>
        <!-- Definitions can have html formatting in them -->
        <span v-html="cardManager.englishDefn" />
      </p>
    </div>
    <div v-if="enableChinese">
      <button
        class="btn-primary btn-xs btn m-2"
        @click="cardManager.changeStep(StepsEnum.CHINESE)"
      >
        Edit
      </button>
      <h2 class="text-l inline font-bold">
        ChineseDefinition:
      </h2>
      <p>
        <!-- Definitions can have html formatting in them -->
        <span v-html="cardManager.chineseDefn" />
      </p>
    </div>
    <div>
      <button
        class="btn-primary btn-xs btn m-2"
        @click="cardManager.changeStep(StepsEnum.IMAGE)"
      >
        Edit
      </button>
      <h2 class="text-l inline font-bold">
        Images:
      </h2>
      <div class="flex gap-1">
        <img
          v-for="image in cardManager.imageUrls"
          :key="image"
          class="max-h-24 w-auto"
          :src="image"
          alt="image for word"
        >
        <img
          v-for="imagedata in cardManager.image64"
          :key="imagedata"
          class="max-h-24 w-auto"
          :src="`data:image/png;base64, ${imagedata}`"
          alt="image for word"
        >
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import { useCardManager } from '@/stores/CardManager';
import { getUserSettings } from '@/lib/userSettings';

const UserSettings = getUserSettings();

const enableChinese = UserSettings.Dictionaries.EnableChinese.read();
const cardManager = useCardManager();

defineEmits(['change-step']);

</script>
