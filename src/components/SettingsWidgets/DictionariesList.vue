<template>
  <n-list class="col-span-2 row-span-4">
    <n-list-item
      v-for="dict in dicts"
      :key="dict">
      {{dict}}
      <template #suffix>
        <n-button @click="deleteDict(dict)">
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

<script setup>
import {
  NModal, NButton, NList, NListItem, NInput, NSelect,
} from 'naive-ui';
import { onBeforeMount, ref } from 'vue';

const dicts = ref([]);
const addDictModal = ref(false);

function deleteDict(dict) {
  console.log('pretend to delete ', dict);
}

const newDictFile = ref('');
const newDictName = ref('');
const newDictType = ref('');
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
  newDictFile.value = await window.ipc.filePicker('json');
}
function addDictionary() {
  // Reset whatever state
  newDictFile.value = '';
  newDictName.value = '';
  newDictType.value = '';
  addDictModal.value = true;
}

function submit() {
  // TODO verify
  window.ipc.addDictionary(
    newDictName.value,
    newDictFile.value,
    newDictType.value,
  );
}

onBeforeMount(async () => {
  dicts.value = await window.ipc.dictionaryInfo();
  console.log(dicts.value);
});

</script>
