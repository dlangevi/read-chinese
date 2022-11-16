<template>
  <div class="container mx-auto px-4">
    <h2 class="mt-5">
      Settings
    </h2>
    <p>Lets go</p>
    <div
      v-for="(contents, section, index) in sections"
      :key="index"
      class="m-4 grid grid-cols-4 justify-start gap-4 border-2 p-4"
    >
      <div class="border-4">
        {{ section }}
      </div>
      <div
        class="col-span-3
               grid
               grid-flow-row
               grid-cols-3
               justify-start
               gap-6"
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
