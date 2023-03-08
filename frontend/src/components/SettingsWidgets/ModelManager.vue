<template>
  <div class="col-span-3 row-span-2 flex flex-col gap-4">
    <div class="flex flex-wrap gap-4">
      <h2 class="text-2xl">
        Configure Anki mappings:
      </h2>
      <select v-model="currentModel" class="select-primary select">
        <option
          v-for="model in models"
          :key="model"
          :value="model"
        >
          {{ model }}
        </option>
      </select>
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
  GetUserSettings,
} from '@wailsjs/backend/UserSettings';
import {
  backend,
} from '@wailsjs/models';

async function generateModel() {
  // call go backend
  await LoadTemplate();
}

type FieldsMapKey = keyof backend.FieldsMapping
const fields = ref(backend.FieldsMapping.createFrom());
const filteredFields = computed((): FieldsMapKey[] => {
  return Object.keys(fields.value).filter((key) => {
    return key !== 'firstField';
  }).map((key) => {
    return key as FieldsMapKey;
    // return key as keyof backend.FieldsMapping;
  });
});
const models = ref<string[]>([]);
const modelFields = ref<string[]>([]);
const currentModel = ref('');
const containsDefault = computed(
  () => models.value.includes('read-chinese-note'),
);

watch(fields, () => {
  console.log(fields.value);
});

watch(currentModel, async () => {
  if (currentModel.value !== '') {
    fields.value = await GetMapping(currentModel.value);
    modelFields.value = await LoadModelFields(currentModel.value);
    fields.value.firstField = modelFields.value[0];
  }
});

function saveMapping() {
  console.log(currentModel.value, fields.value);
  SetMapping(currentModel.value, fields.value);
}

async function fetchModels() {
  return LoadModels().then(async (loaded) => {
    if (loaded !== undefined) {
      models.value = loaded;

      // Ideally we would pass the usersettings in somehow
      // But this works for now
      const userSettings = await GetUserSettings();
      const activeModel = userSettings.AnkiConfig.ActiveModel;
      if (loaded.includes(activeModel)) {
        currentModel.value = activeModel;
      }
    } else {
      console.log('Failed to fetch data for selector');
      setTimeout(fetchModels, 1000);
    }
  }).catch((err) => {
    console.log(err);
    setTimeout(fetchModels, 1000);
  });
}

onBeforeMount(async () => {
  await fetchModels();
});
</script>
