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
      <div class="m-auto flex w-3/4 flex-col gap-1">
        <n-input
          v-model:value="newDictName"
          type="text"
          placeholder="Dictionary Name"
        />
        <div>
          <n-button @click="pickFile">
            Select File
          </n-button>
          {{ newDictFile }}
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
      <template #header>
        <div class="grid grid-cols-2 text-2xl">
          <div>
            Dictionaries
          </div>
          <n-button @click="addDictionary">
            Add Dictionary
          </n-button>
        </div>
      </template>
      <n-list-item
        v-for="(dict, name) in dicts"
        :key="name"
      >
        <n-descriptions label-placement="left">
          <n-descriptions-item label="Name">
            {{ dict.name }}
          </n-descriptions-item>
          <n-descriptions-item label="Type">
            {{ dict.type }}
          </n-descriptions-item>
          <n-descriptions-item label="Path">
            {{ dict.path }}
          </n-descriptions-item>
        </n-descriptions>
        <template #suffix>
          <n-button @click="makePrimary(name as string)">
            Make Primary
          </n-button>
          <n-button @click="deleteDict(name as string)">
            Delete
          </n-button>
        </template>
      </n-list-item>
    </n-list>
  </div>
</template>

<script lang="ts" setup>
import {
  NModal, NButton, NList, NListItem, NInput, NSelect,
  NDescriptions, NDescriptionsItem,
} from 'naive-ui';
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
