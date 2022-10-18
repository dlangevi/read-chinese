<template>
  <n-list class="col-span-2 row-span-4">
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
</template>

<script lang="ts" setup>
import {
  NModal, NButton, NList, NListItem, NInput, NSelect,
} from 'naive-ui';
import { onBeforeMount, ref } from 'vue';
import type { DictionaryInfo, DictionaryType } from '@/shared/types';

const addDictModal = ref(false);

function makePrimary(name:string) {
  window.nodeIpc.setPrimaryDict(name);
}

function deleteDict(name:string) {
  window.nodeIpc.deleteDictionary(name);
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
  newDictFile.value = await window.nodeIpc.filePicker('json');
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
  window.nodeIpc.addDictionary(
    newDictName.value,
    newDictFile.value,
    newDictType.value,
  );
}

const dicts = ref<{ [name:string]: DictionaryInfo }>({});
onBeforeMount(async () => {
  dicts.value = await window.nodeIpc.dictionaryInfo();
  console.log(dicts.value);
});

</script>
