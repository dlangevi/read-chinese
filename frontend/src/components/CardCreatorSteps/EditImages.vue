<template>
  <div>
    <div class="m-4 text-3xl">
      Pick an image
    </div>
    <div class="flex flex-wrap">
      <div
        v-for="(imageData, i) in images"
        :key="i"
        class="max-w-[33%]"
      >
        <label
          class="
          label
          m-4
          cursor-pointer gap-2"
          :for="i.toString()"
        >

          <input
            :id="i.toString()"
            v-model="selectedImages"
            class="checkbox"
            :value="imageData"
            type="checkbox"
            name="images"
          >
          <img
            class="h-auto w-auto"
            :src="imageData.url"
            :alt="imageData.name ||
              'Image related to search word, no alt text generated'"
          >
        </label>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  watch, ref,
} from 'vue';
import { SearchImages } from '@wailsjs/backend/ImageClient';
import { backend } from '@wailsjs/models';
import { useCardManager } from '@/stores/CardManager';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import { storeToRefs } from 'pinia';
import { getUserSettings } from '@/lib/userSettings';

const UserSettings = getUserSettings();
const cardManager = useCardManager();
const { currentStep } = storeToRefs(cardManager);

const images = ref<backend.ImageInfo[]>([]);
const selectedImages = ref<backend.ImageInfo[]>([]);

watch(selectedImages, async () => {
  cardManager.updateImages(selectedImages.value);
  const autoAdvance = UserSettings.CardCreation.AutoAdvanceImage;
  if (autoAdvance) {
    cardManager.nextStep();
  }
});

watch(currentStep, async () => {
  if (currentStep.value === StepsEnum.IMAGE && images.value.length === 0) {
    images.value = await SearchImages(cardManager.word);
    console.log('images', images.value);
  }
});
</script>
