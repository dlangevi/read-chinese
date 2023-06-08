<template>
  <div class="col-span-3 row-span-2 flex flex-col gap-2">
    <div
      :class="['modal', {'modal-open': addDictModal}]"
      @click="() => addDictModal = false"
    >
      <div
        class="modal-box relative w-1/2 "
        @click.stop
      >
        <div>
          Add a dictionary in Migaku Json Format
        </div>
        <div class="divider" />
        <div class="m-auto flex w-3/4 flex-col gap-1">
          <input
            v-model="newDictName"
            class="input-bordered input"
            type="text"
            placeholder="Dictionary Name"
          >
          <div>
            <button class="btn-secondary btn" @click="pickFile">
              Select File
            </button>
            {{ newDictFile }}
          </div>
          <select
            v-model="newDictType"
            class="select-bordered select"
          >
            <option
              v-for="(option, i) in options"
              :key="i"
              :value="option.value"
            >
              {{ option.label }}
            </option>
          </select>
        </div>
        <div class="modal-action">
          <button class="btn-secondary btn" @click="submit">
            Add Custom Dictionary
          </button>
        </div>
      </div>
    </div>
    <div class="flex items-center gap-4">
      <div class="grow text-2xl">
        Dictionaries
      </div>
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <button class="btn-secondary btn justify-end" @click="addCedict">
            Add Default Dict
          </button>
          <div
            class="tooltip tooltip-left"
            data-tip="Downloads CC-Cedict"
          >
            <information-circle class="h-6 w-6" />
          </div>
        </div>
        <div class="flex items-center gap-2">
          <button class="btn-secondary btn" @click="addDictionary">
            Add Custom Dictionary
          </button>
          <div
            class="tooltip tooltip-left"
            data-tip="Import your own json dictionary"
          >
            <information-circle class="h-6 w-6" />
          </div>
        </div>
      </div>
    </div>
    <table class="table">
      <thead>
        <tr>
          <th>Primary</th>
          <th>Name</th>
          <th>Type</th>
          <th>Path</th>
          <th />
          <th />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(dict, name) in dicts"
          :key="name"
        >
          <td class="">
            <star
              v-if="dict.isPrimary"
              class="inline h-6 w-6"
            />
          </td>
          <td class="">
            {{ dict.name }}
          </td>
          <td> {{ dict.type }} </td>
          <td> {{ dict.path }} </td>
          <td>
            <button
              class="btn-secondary btn"
              @click="makePrimary(name as string)"
            >
              Make Primary
            </button>
          </td>
          <td>
            <button
              class="btn-secondary btn"
              @click="deleteDict(name as string)"
            >
              Delete
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref } from 'vue';
import { InformationCircle, Star } from '@vicons/ionicons5';

import {
  GetDictionaryInfo,
  SetPrimaryDict,
  DeleteDictionary,
  AddCedict,
  AddMigakuDictionary,
} from '@wailsjs/backend/Dictionaries';
import {
  backend,
} from '@wailsjs/models';

import { FilePicker } from '@wailsjs/backend/Backend';

import { useLoader } from '@/lib/loading';

// TODO see if the go type system can represent these
export type DictionaryType = 'english' | 'chinese';

const addDictModal = ref(false);

async function makePrimary(name:string) {
  await SetPrimaryDict(name);
  updateDicts();
}

async function deleteDict(name:string) {
  await DeleteDictionary(name);
  updateDicts();
}

const newDictFile = ref('');
const newDictName = ref('');
const newDictType = ref<DictionaryType>('english');
const options = [
  {
    label: 'English - Chinese',
    value: 'english',
  },
  {
    label: 'Chinese - Chinese',
    value: 'chinese',
  },
];

async function pickFile() {
  newDictFile.value = await FilePicker('json');
}
const loader = useLoader();
async function addCedict() {
  await loader.withLoader(AddCedict);
  updateDicts();
}

function addDictionary() {
  // Reset whatever state
  newDictFile.value = '';
  newDictName.value = '';
  newDictType.value = 'english';
  addDictModal.value = true;
}

async function submit() {
  // TODO verify
  await AddMigakuDictionary(
    newDictName.value,
    newDictFile.value,
    newDictType.value,
  );
  updateDicts();
  addDictModal.value = false;
}

const dicts = ref<{[key:string]:backend.DictionaryInfo}>({});
async function updateDicts() {
  dicts.value = await GetDictionaryInfo();
}
onBeforeMount(async () => {
  updateDicts();
});

</script>
