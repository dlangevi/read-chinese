<template>
  <div class="flex flex-col gap-4 p-4">
    <div
      v-for="(contents, section, index) in sections"
      :key="index"
      class="
      mx-auto
      flex
      w-4/5
      border-2
      border-base-300
      p-4"
    >
      <div class="m-4 w-1/6 text-2xl font-extrabold">
        {{ section }}
      </div>
      <div
        class="grid
               w-5/6
               grid-flow-row
               grid-cols-3
               gap-4"
      >
        <component
          :is="content.type"
          v-for="([name, content]) in contents"
          :key="name"
          :class="[{'bg-sky-500': props.highlight === content.name}]"
          :setting="content"
          :initial-value="getInitialValue(section, name)"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  ComponentTable, getUserSettings,
} from '@/lib/userSettings';
import { backend } from '@wailsjs/models';

const props = defineProps<{
  highlight?: string,
}>();

const UserSettings = getUserSettings();

function getInitialValue(section: string, value: string) {
  const segment = UserSettings[section as keyof backend.UserSettings];
  const segvalue : unknown = segment[value as keyof typeof segment];
  return segvalue;
}

// TODO the way the initial value is looked up is not goo
const sections = {
  meta:
    Object.entries(ComponentTable.meta),
  CardCreation:
     Object.entries(ComponentTable.CardCreation),
  AnkiConfig:
    Object.entries(ComponentTable.AnkiConfig),
  AzureConfig:
    Object.entries(ComponentTable.AzureConfig),
  Dictionaries:
     Object.entries(ComponentTable.Dictionaries),
  WordLists:
     Object.entries(ComponentTable.WordLists),
  SentenceGeneration:
    Object.entries(ComponentTable.SentenceGeneration),
};
</script>
