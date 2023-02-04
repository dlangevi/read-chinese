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
            v-model="image"
            class="checkbox"
            :value="i"
            type="checkbox"
            name="images"
          >
          <img
            class="h-auto w-auto"
            :src="imageData.thumbnailUrl"
            alt="image for word"
          >
        </label>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  watch, onBeforeMount, ref,
} from 'vue';
import { SearchImages } from '@wailsjs/backend/ImageClient';
import { backend } from '@wailsjs/models';
import { useCardManager } from '@/stores/CardManager';
import { storeToRefs } from 'pinia';
import { getUserSettings } from '@/lib/userSettings';

const UserSettings = getUserSettings();
const cardManager = useCardManager();
const { word } = storeToRefs(cardManager);

const images = ref<backend.ImageInfo[]>([]);
const image = ref([]);

watch(image, async () => {
  const entries = image.value.map((i) => images.value[i]);
  cardManager.updateImages(entries);
  const autoAdvance = UserSettings.CardCreation.AutoAdvanceImage.read();
  if (autoAdvance) {
    cardManager.nextStep();
  }
});

watch(word, () => {
  loadData();
});

async function loadData() {
  images.value = await SearchImages(cardManager.word);
}

onBeforeMount(loadData);
</script>
