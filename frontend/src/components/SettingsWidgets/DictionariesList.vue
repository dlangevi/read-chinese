<template>
  <div class="col-span-3 row-span-2">
    <n-modal
      v-model:show="addDictModal"
      :closable="true"
      class="w-1/2"
      preset="card"
    >
      <template #header>
        Add a dicitonary in Yomichan Json Format
      </template>
      <div class="w-3/4 m-auto flex flex-col gap-1">
        <n-input
          v-model:value="newDictName"
          type="text"
          placeholder="Dictionary Name" />
        <div>
          <n-button @click="pickFile">
            Select File
          </n-button>
          {{newDictFile}}
        </div>
        <n-select
          v-model:value="newDictType"
          :options="options"
        />
      </div>
      <template #action>
        <n-button @click="submit">
          Add Dictionary
        </n-button>
      </template>
    </n-modal>
    <n-list class="col-span-3 row-span-2">
      <n-list-item
        v-for="(dict, name) in dicts"
        :key="name">
        {{dict}}
        <template #suffix>
          <n-button @click="makePrimary(name as string)">
            Make Primary
          </n-button>
          <n-button @click="deleteDict(name as string)">
            Delete
          </n-button>
        </template>
      </n-list-item>
      <template #footer>
        <n-button @click="addDictionary">
          Add Dictionary
        </n-button>
      </template>
    </n-list>
  </div>
</template>

<script lang="ts" setup>
import {
  NModal, NButton, NList, NListItem, NInput, NSelect,
} from 'naive-ui';
import { onBeforeMount, ref } from 'vue';
import type { DictionaryInfo, DictionaryType } from '@/lib/types';

import {
  GetDictionaryInfo,
  SetPrimaryDict,
  DeleteDictionary,
  AddDictionary,
} from '@wailsjs/backend/Dictionaries';

import { FilePicker } from '@wailsjs/main/App';

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
}

const dicts = ref<{ [name:string]: DictionaryInfo }>({});
onBeforeMount(async () => {
  dicts.value = await GetDictionaryInfo();
});

</script>
