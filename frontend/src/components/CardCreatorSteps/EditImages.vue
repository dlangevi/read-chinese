<template>
  <div>
    <div class="m-4 text-3xl">
      Pick an image
    </div>
    <div class="flex">
      <!-- rows -->
      <div
        v-for="(col, c) in [cols1, cols2, cols3]"
        :key="c"
        class="flex flex-wrap gap-2"
      >
        <!-- column -->
        <div>
          <label
            v-for="(image, i) in col"
            :key="i"
            :for="`${c*100 + i}`"
            class="label w-full cursor-pointer"
          >
            <div class="flex w-full place-items-center gap-2">
              <input
                :id="`${c*100 + i}`"
                v-model="selectedImages"
                class="checkbox"
                :value="image"
                type="checkbox"
                name="images"
              >
              <img
                class="w-1/2 flex-initial grow"
                :src="getImageSrc(image)"
                :alt="image.name ||
                  'Image related to search word, no alt text generated'"
              >
            </div>
          </label>
        </div>
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
const cols1 = ref<backend.ImageInfo[]>([]);
const cols2 = ref<backend.ImageInfo[]>([]);
const cols3 = ref<backend.ImageInfo[]>([]);

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
    images.value.push(...await SearchImages(cardManager.word));

    // Sort items in descending order based on height with
    // original selected images in first
    images.value.sort((a, b) => {
      if (selectedImages.value.includes(a)) {
        return -1;
      }
      if (selectedImages.value.includes(b)) {
        return 1;
      }
      return b.imageHeight - a.imageHeight;
    });

    // Create three empty sublists
    const A: backend.ImageInfo[] = [];
    const B: backend.ImageInfo[] = [];
    const C: backend.ImageInfo[] = [];

    // Initialize total heights for each sublist
    const totalHeights = [0, 0, 0];

    // Iterate over each item in the sorted list
    for (const item of images.value) {
      // Find the sublist with the smallest total height so far
      const minIndex = totalHeights.indexOf(Math.min(...totalHeights));

      // Add the item to the selected sublist
      if (minIndex === 0) {
        A.push(item);
      } else if (minIndex === 1) {
        B.push(item);
      } else {
        C.push(item);
      }

      // Update the total height of the selected sublist
      totalHeights[minIndex] += item.imageHeight;
    }

    // Randomize the individual lists, with already selected images
    // up first (because having all long images first looks weird)
    const sorter = (a : backend.ImageInfo, b : backend.ImageInfo) => {
      if (selectedImages.value.includes(a)) {
        return -1;
      }
      if (selectedImages.value.includes(b)) {
        return 1;
      }
      return Math.random() - 0.5;
    };
    A.sort(sorter);
    B.sort(sorter);
    C.sort(sorter);
    // Return the three sublists
    cols1.value = A;
    cols2.value = B;
    cols3.value = C;
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
