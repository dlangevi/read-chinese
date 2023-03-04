<template>
  <div class="mx-auto px-4">
    <div
      v-for="(contents, section, index) in sections"
      :key="index"
      class="m-4 grid grid-cols-4 justify-start
      gap-4 border-2 border-base-300 p-4"
    >
      <div class="m-4 text-2xl font-extrabold">
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
import {
  getUserSettings, getDisplayable,
} from '@/lib/userSettings';

const props = defineProps<{
  highlight?: string,
}>();

const UserSettings = getUserSettings();

const sections = {
  GeneralSettings:
    getDisplayable(UserSettings.meta),
  CardCreationSettings:
    getDisplayable(UserSettings.CardCreation),
  AnkiSettings:
    getDisplayable(UserSettings.AnkiConfig),
  DictionarySettings:
    getDisplayable(UserSettings.Dictionaries),
  SentenceGeneration:
    getDisplayable(UserSettings.SentenceGeneration),
};
</script>
