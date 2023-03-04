<template>
  <div>
    <div class="m-4 text-3xl">
      Pick an image
    </div>
    <div class="flex flex-wrap">
      <div
        v-for="(image, i) in images"
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
            :value="image"
            type="checkbox"
            name="images"
          >
          <img
            class="h-auto w-auto"
            :src="getImageSrc(image)"
            :alt="image.name ||
              'Image related to search word, no alt text generated'"
          >
        </label>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  onBeforeMount, watch, ref, toRaw,
} from 'vue';
import { SearchImages } from '@wailsjs/backend/ImageClient';
import { backend } from '@wailsjs/models';
import { useCardManager, getImageSrc } from '@/stores/CardManager';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import { getUserSettings } from '@/lib/userSettings';
let allowUpdate = false;
let loaded = false;

const UserSettings = getUserSettings();
const cardManager = useCardManager();

const images = ref<backend.ImageInfo[]>([]);
const selectedImages = ref<backend.ImageInfo[]>([]);

watch(selectedImages, async () => {
  if (!allowUpdate) {
    return;
  }
  cardManager.updateImages(
    selectedImages.value.map((img) => toRaw(img)));
  const autoAdvance = UserSettings.CardCreation.AutoAdvanceImage;
  if (autoAdvance) {
    cardManager.nextStep();
  }
});

watch(() => cardManager.currentStep, async () => {
  if (cardManager.currentStep === StepsEnum.IMAGE && !loaded) {
    loaded = true;
    console.log('doing image query now');
    images.value.push(...await SearchImages(cardManager.word));
  }

  allowUpdate = true;
}, {
  immediate: true,
});

onBeforeMount(async () => {
  const originalImages = cardManager.originalValues?.images;
  if (originalImages) {
    images.value.push(...originalImages);
    selectedImages.value.push(...originalImages);
  }
});

</script>
