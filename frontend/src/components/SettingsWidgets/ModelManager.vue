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
        v-for="(_, key) in fields"
        :key="key"
        class="flex place-content-between"
      >
        <label>{{ key }}</label>
        <select
          v-model="fields[key]"
          class="select-primary select"
        >
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
} from '@wailsjs/backend/UserConfig';
import {
  backend,
} from '@wailsjs/models';

async function generateModel() {
  // call go backend
  await LoadTemplate();
}

const fields = ref(backend.FieldsMapping.createFrom());
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
  }
});

function saveMapping() {
  console.log(currentModel.value, fields.value);
  SetMapping(currentModel.value, fields.value);
}

onBeforeMount(async () => {
  models.value = await LoadModels();

  // Ideally we would pass the usersettings in somehow
  // But this works for now
  const userSettings = await GetUserSettings();
  let activeModel = userSettings.AnkiConfig.ActiveModel;
  fields.value = await GetMapping(activeModel);
  if (models.value.length > 0) {
    if (!models.value.includes(activeModel)) {
      activeModel = models.value[0];
    }
    modelFields.value = await LoadModelFields(activeModel);
  }
});
</script>
