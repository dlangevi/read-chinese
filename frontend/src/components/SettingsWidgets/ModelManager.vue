<template>
  <div class="col-span-3 row-span-2 flex flex-wrap gap-4">
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
  </div>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref, watch } from 'vue';
import {
  LoadTemplate,
  LoadModels,
  LoadModelFields,
} from '@wailsjs/backend/ankiInterface';
import {
  GetMapping,
  GetUserSetting,
} from '@wailsjs/backend/UserSettings';
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

watch(fields, () => {
  console.log(fields.value);
});

onBeforeMount(async () => {
  models.value = await LoadModels();
  const activeModel = await GetUserSetting('ActiveModel');
  modelFields.value = await LoadModelFields(activeModel);
  console.log(fields);
  fields.value = await GetMapping(activeModel);
});
</script>
