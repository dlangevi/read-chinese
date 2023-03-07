<template>
  <div class="mx-auto w-4/5 px-4">
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
  const segment = UserSettings[section as keyof backend.UserConfig];
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
  Dictionaries:
     Object.entries(ComponentTable.Dictionaries),
  SentenceGeneration:
    Object.entries(ComponentTable.SentenceGeneration),
};
</script>
