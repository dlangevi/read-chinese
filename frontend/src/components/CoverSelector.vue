<template>
  <div>
    <button class="btn-secondary btn" @click="getCovers">
      Select Cover Image from ImageSearch
    </button>
    <div class="flex flex-wrap">
      <label
        v-for="(image, i) in images"
        :key="i"
        :for="image.url"
        class="label w-1/5 cursor-pointer"
      >
        <div class="flex  place-items-center gap-2">
          <input
            :id="image.url"
            v-model="selectedImage"
            class="radio"
            :value="image"
            type="radio"
            name="covers"
          >
          <img
            class="w-1/2 flex-initial grow"
            :src="image.url"
            :alt="image.name ||
              'Image related to search word, no alt text generated'"
          >
        </div>
      </label>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { backend } from '@wailsjs/models';
import { SearchBookCovers } from '@wailsjs/backend/ImageClient';
const images = ref<backend.ImageInfo[]>([]);
const selectedImage = ref<backend.ImageInfo>();
async function getCovers() {
  // TODO get from props
  images.value = await SearchBookCovers('', '');
  // bookAuthor.value, bookTitle.value);
}

</script>
