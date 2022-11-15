<template>
  <div>
    <div class="text-3xl m-4">
      Pick an image
    </div>
    <n-checkbox-group
      v-model:value="image"
      name="images"
    >
      <n-space :size="40">
        <n-checkbox
          v-for="(imageData, i) in images"
          :key="i"
          :value="i"
        >
          <img
            class="h-48 w-auto"
            :src="imageData.thumbnailUrl"
            alt="image for word"
          >
        </n-checkbox>
      </n-space>
    </n-checkbox-group>
  </div>
</template>

<script lang="ts" setup>
import {
  watch, onBeforeMount, ref,
} from 'vue';
import {
  NSpace, NCheckboxGroup, NCheckbox,
} from 'naive-ui';
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
