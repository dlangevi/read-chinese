<template>
  <div class="container mx-auto px-4">
    <h2 class="mt-5">Settings</h2>
    <p>Lets go</p>
    <div
      v-for="(contents, section, index) in sections"
      :key="index"
      class="grid gap-4 grid-cols-4 m-4 p-4 border-2 justify-start">
      <div class="border-4">
        {{section}}
      </div>
      <div
        class="col-span-3 grid gap-4 grid-cols-3 justify-start">
        <component
          class="h-12"
          v-for="content in contents"
          :key="content"
          :is="content.type"
          :setting="content"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue';

const UserSettings = inject('userSettings');

const sections = {
  CardCreationSettings:
    [...Object.values(UserSettings.CardCreation)],
  DictionarySettings:
    Object.values(UserSettings.Dictionaries),

};
</script>
