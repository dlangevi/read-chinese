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
import { getUserSettings } from '@/lib/userSettings';
import { SearchImages } from '@wailsjs/backend/ImageClient';
import { backend } from '@wailsjs/models';

const UserSettings = getUserSettings();

const emit = defineEmits(['update-images']);
const images = ref<backend.ImageInfo[]>([]);
const image = ref([]);
const loaded = ref(false);

watch(image, async () => {
  const autoAdvance = UserSettings.CardCreation.AutoAdvanceImage.read();
  const entries = image.value.map((i) => images.value[i]);
  emit('update-images', entries, autoAdvance);
});

const props = defineProps<{
  word: string,
}>();

onBeforeMount(async () => {
  images.value = await SearchImages(props.word);
  loaded.value = true;
});

</script>
