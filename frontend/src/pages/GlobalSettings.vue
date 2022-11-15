<template>
  <div class="container mx-auto px-4">
    <h2 class="mt-5">
      Settings
    </h2>
    <p>Lets go</p>
    <div
      v-for="(contents, section, index) in sections"
      :key="index"
      class="grid gap-4 grid-cols-4 m-4 p-4 border-2 justify-start"
    >
      <div class="border-4">
        {{ section }}
      </div>
      <div
        class="col-span-3 grid gap-6 grid-cols-3 grid-flow-row justify-start"
      >
        <component
          :is="content.type"
          v-for="content in contents"
          :key="content"
          :setting="content"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { getUserSettings } from '@/lib/userSettings';

const UserSettings = getUserSettings();

const sections = {
  CardCreationSettings:
    [...Object.values(UserSettings.CardCreation)],
  DictionarySettings:
    Object.values(UserSettings.Dictionaries),

};
</script>
