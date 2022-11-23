<template>
  <div class="col-span-3 row-span-2">
    <div
      :class="['modal', {'modal-open': addDictModal}]"
      @click="() => addDictModal = false"
    >
      <div class="modal-box relative w-1/2 max-w-5xl">
        <div>
          Add a dicitonary in Yomichan Json Format
        </div>
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
          <button class="btn-secondary btn" @click="submit">
            Add Dictionary
          </button>
        </div>
      </div>
    </div>
    <div class="grid grid-cols-2 text-2xl">
      <div>
        Dictionaries
      </div>
      <button class="btn-secondary btn" @click="addDictionary">
        Add Dictionary
      </button>
    </div>
    <div class="col-span-3 row-span-2">
      <div
        v-for="(dict, name) in dicts"
        :key="name"
      >
        <div class="divider" />
        <div class="grid grid-cols-4">
          <div class="col-span-3 flex gap-3">
            <div> Name: {{ dict.name }} </div>
            <div> Type: {{ dict.type }} </div>
            <div> Path: {{ dict.path }} </div>
          </div>
          <div class="grid gap-4">
            <button
              class="btn-secondary btn"
              @click="makePrimary(name as string)"
            >
              Make Primary
            </button>
            <button
              class="btn-secondary btn"
              @click="deleteDict(name as string)"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref } from 'vue';

import {
  GetDictionaryInfo,
  SetPrimaryDict,
  DeleteDictionary,
  AddDictionary,
} from '@wailsjs/backend/Dictionaries';

import { FilePicker } from '@wailsjs/main/App';

// TODO see if the go type system can represent these
export type DictionaryType = 'english' | 'chinese';

const addDictModal = ref(false);

function makePrimary(name:string) {
  SetPrimaryDict(name);
}

function deleteDict(name:string) {
  DeleteDictionary(name);
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
function addDictionary() {
  // Reset whatever state
  newDictFile.value = '';
  newDictName.value = '';
  newDictType.value = 'english';
  addDictModal.value = true;
}

function submit() {
  // TODO verify
  AddDictionary(
    newDictName.value,
    newDictFile.value,
    newDictType.value,
  );
  addDictModal.value = false;
}

// TODO backend does not export these types
// even though IMO it should
type DictionaryInfo = {
  name: string,
  path: string,
  type: string,
};
type DictionaryInfoMap = {
  [name:string] : DictionaryInfo
};

const dicts = ref<DictionaryInfoMap>({});
onBeforeMount(async () => {
  dicts.value = await GetDictionaryInfo();
});

</script>
