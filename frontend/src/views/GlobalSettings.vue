<template>
  <div class="container mx-auto px-4">
    <h2 class="mt-5">
      Settings
    </h2>
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
        <div class="m-4 text-2xl font-extrabold">
          {{ section }}
        </div>
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
          v-for="([initial, content]) in contents"
          :key="content.name"
          :class="[{'bg-sky-500': props.highlight === content.name}]"
          :setting="content"
          :initial-value="initial"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { getUserSettings, getDisplayable } from '@/lib/userSettings';
import { onMounted } from 'vue';
import { themeChange } from 'theme-change';

const props = defineProps<{
  highlight?: string,
}>();
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
    getDisplayable(UserSettings.CardCreation),
  AnkiSettings:
    getDisplayable(UserSettings.AnkiConfig),
  DictionarySettings:
    getDisplayable(UserSettings.Dictionaries),
  SentenceGeneration:
    getDisplayable(UserSettings.SentenceGeneration),
};
console.log(UserSettings.Dictionaries);
console.log(sections);
</script>
