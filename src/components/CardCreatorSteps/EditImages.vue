<template>
  <div class="text-3xl m-4">Pick an image</div>
  <n-checkbox-group v-model:value="image" name="images">
    <n-space :size="40">
        <n-checkbox
          v-for="(imageData, i) in images"
          :key="i"
          :value="imageData"
        >
        <img class="h-48 w-auto"
        :src="imageData.thumbnailUrl" alt="image for word"/>
      </n-checkbox>
    </n-space>
  </n-checkbox-group>
</template>

<script setup>
import {
  watch, onBeforeMount, ref, inject,
} from 'vue';
import {
  NSpace, NCheckboxGroup, NCheckbox,
} from 'naive-ui';

const UserSettings = inject('userSettings');

const emit = defineEmits(['updateImages']);
const images = ref([]);
const image = ref(null);
const loaded = ref(false);

watch(image, async () => {
  const autoAdvance = UserSettings.CardCreation.AutoAdvanceImage.read();
  emit('updateImages', image.value, autoAdvance);
});

const props = defineProps({
  word: {
    type: String,
    required: true,
  },
});

onBeforeMount(async () => {
  images.value = await window.ipc.getImagesForWord(props.word);
  loaded.value = true;
});

</script>
