<template>
  <div class="col-span-3 row-span-2 flex flex-col gap-2">
    <div
      :class="['modal', {'modal-open': addListModal}]"
      @click="() => addListModal = false"
    >
      <div
        class="modal-box relative w-5/6"
        @click.stop
      >
        <div>
          Add frequency list
        </div>
        <div class="divider" />
        <div class="m-auto flex flex-col gap-1">
          <div class="flex items-center">
            <button class="btn-secondary btn" @click="pickFile">
              Select File
            </button>
            <div>
              {{ newListFile }}
            </div>
          </div>
          <input
            v-model="newListName"
            class="input-bordered input"
            type="text"
            placeholder="List Name"
          >
        </div>
        <div class="modal-action">
          <button class="btn-secondary btn" @click="submit">
            Add Frequency List
          </button>
        </div>
      </div>
    </div>
    <div class="flex items-center gap-4">
      <div class="grow text-2xl">
        Word Lists
      </div>
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <button class="btn-secondary btn" @click="addWordList">
            Add Word List
          </button>
          <div
            class="tooltip tooltip-left"
            data-tip="Import your own frequency list"
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
          <th>Path</th>
          <th />
          <th />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(list, name) in lists"
          :key="name"
        >
          <td class="">
            <star
              v-if="list.isPrimary"
              class="inline h-6 w-6"
            />
          </td>
          <td class="">
            {{ list.name }}
          </td>
          <td> {{ list.path }} </td>
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
              @click="deleteList(name as string)"
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
  SetPrimaryList,
  DeleteList,
  AddList,
  GetWordListsInfo,
} from '@wailsjs/backend/wordLists';
import {
  backend,
} from '@wailsjs/models';

import { FilePicker, GetFileName } from '@wailsjs/backend/Backend';

const addListModal = ref(false);

async function makePrimary(name:string) {
  await SetPrimaryList(name);
  updateLists();
}

async function deleteList(name:string) {
  await DeleteList(name);
  updateLists();
}

const newListFile = ref('');
const newListName = ref('');

async function pickFile() {
  newListFile.value = await FilePicker('txt');
  if (newListName.value === '') {
    const fileName = await GetFileName(newListFile.value);
    newListName.value = fileName;
  }
}

function addWordList() {
  // Reset whatever state
  newListFile.value = '';
  newListName.value = '';
  addListModal.value = true;
}

async function submit() {
  // TODO verify
  await AddList(
    newListName.value,
    newListFile.value,
  );
  updateLists();
  addListModal.value = false;
}

const lists = ref<{[key:string]:backend.WordListInfo}>({});
async function updateLists() {
  lists.value = await GetWordListsInfo();
}
onBeforeMount(async () => {
  updateLists();
});

</script>
