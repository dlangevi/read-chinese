<template>
  <div class="col-span-3 row-span-2 flex flex-col gap-4">
    <div class="flex flex-wrap gap-4">
      <h2 class="text-2xl">
        Anki mappings for: {{ UserSettings.AnkiConfig.ActiveModel }}
      </h2>
      <div class="btn-primary btn" @click="saveMapping">
        Save current mapping
      </div>
      <div
        v-if="!containsDefault"
        class="btn-primary btn"
        @click="generateModel"
      >
        Install Default
      </div>
    </div>
    <div class="grid grid-cols-3 gap-4">
      <div
        v-for="(key) in filteredFields"
        :key="key"
        class="flex place-content-between"
      >
        <label>{{ key }}</label>
        <select
          v-model="fields[key]"
          class="select-primary select"
        >
          <option value="">
            Unused
          </option>
          <option
            v-for="field in modelFields"
            :key="field"
            :value="field"
          >
            {{ field }}
          </option>
        </select>
      </div>
    </div>
    <div class="divider" />
  </div>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref, watch, computed } from 'vue';
import {
  LoadTemplate,
  LoadModels,
  LoadModelFields,
} from '@wailsjs/backend/ankiInterface';
import {
  SetMapping,
  GetMapping,
} from '@wailsjs/backend/UserSettings';
import {
  backend,
} from '@wailsjs/models';

import { getUserSettings } from '@/lib/userSettings';
const UserSettings = getUserSettings();

async function generateModel() {
  // call go backend
  await LoadTemplate();
  await fetchModels();
}

type FieldsMapKey = keyof backend.FieldsMapping
const fields = ref(backend.FieldsMapping.createFrom());
const filteredFields = computed((): FieldsMapKey[] => {
  return Object.keys(fields.value).filter((key) => {
    return key !== 'firstField';
  }).map((key) => {
    return key as FieldsMapKey;
  });
});
const models = ref<string[]>([]);
const modelFields = ref<string[]>([]);
const currentModel = computed(() => UserSettings.AnkiConfig.ActiveModel);
const containsDefault = computed(
  () => models.value.includes('read-chinese-note'),
);

watch(
  () => UserSettings.AnkiConfig.ActiveModel,
  async (activeModel) => {
    loadModel(activeModel);
  });

async function loadModel(activeModel : string) {
  if (activeModel !== '') {
    fields.value = await GetMapping(activeModel);
    modelFields.value = await LoadModelFields(activeModel);
    fields.value.firstField = modelFields.value[0];
  }
}

function saveMapping() {
  SetMapping(currentModel.value, fields.value);
}

async function fetchModels() {
  return LoadModels().then(async (loaded) => {
    models.value = loaded;
  }).catch((err) => {
    console.log(err);
    setTimeout(fetchModels, 1000);
  });
}

onBeforeMount(async () => {
  await fetchModels();
  loadModel(UserSettings.AnkiConfig.ActiveModel);
});
</script>
