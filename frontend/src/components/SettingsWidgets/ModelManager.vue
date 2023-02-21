<template>
  <div class="col-span-3 row-span-2 flex gap-4">
    <div class="btn-primary btn" @click="generateModel">GenerateModel</div>
    <select class="select-primary select">
      <option
        v-for="model in models"
        :key="model"
        :value="model"
      >
        {{ model }}
      </option>
    </select>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref } from 'vue';
import {
  LoadTemplate,
  LoadModels,
  // LoadModelFields,
} from '@wailsjs/backend/ankiInterface';

async function generateModel() {
  // call go backend
  await LoadTemplate();
}
const models = ref<string[]>([]);

onBeforeMount(async () => {
  models.value = await LoadModels();
  // console.log(await LoadModelFields('Basic'));
});
</script>
