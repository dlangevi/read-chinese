<template>
  <div class="container mx-auto px-4">
    <h2 class="mt-5">
      Settings
    </h2>
    <p>Lets go</p>
    <div>
      <div class="m-5">
        <h2 class="mb-4 text-2xl font-bold text-green-700">Select Theme</h2>
        <select
          data-choose-theme
          class="select-primary select"
        >
          <option
            v-for="theme in themes"
            :key="theme"
            :value="theme"
          >
            {{ theme }}
          </option>
        </select>
      </div>
    </div>
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
import { onMounted } from 'vue';
import { themeChange } from 'theme-change';
onMounted(() => {
  themeChange(false);
});

const UserSettings = getUserSettings();
const themes = [
  'light', 'dark', 'cupcake',
  'bumblebee', 'emerald', 'corporate',
  'synthwave', 'retro', 'cyberpunk',
  'valentine', 'halloween', 'garden',
  'forest', 'aqua', 'lofi', 'pastel',
  'fantasy', 'wireframe', 'black',
  'luxury', 'dracula', 'cmyk',
  'autumn', 'business', 'acid',
  'lemonade', 'night', 'coffee', 'winter',
];

const sections = {
  CardCreationSettings:
    [...Object.values(UserSettings.CardCreation)],
  DictionarySettings:
    Object.values(UserSettings.Dictionaries),

};
</script>
