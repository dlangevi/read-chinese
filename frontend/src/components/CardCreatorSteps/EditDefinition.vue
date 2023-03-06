<template>
  <div>
    <div class="m-4 text-3xl">
      Pick a definition
    </div>
    <div
      v-for="(def, i) in definitions"
      :key="i"
      class="text-3xl"
    >
      <label
        class="label cursor-pointer justify-start gap-2"
        :for="def.definition + i"
      >
        <input
          :id="def.definition + i"
          v-model="selectedDefs"
          class="checkbox"
          :value="def"
          type="checkbox"
          name="definitions"
        >
        {{ formatDefinition(def) }}
      </label>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  watch, onBeforeMount, ref,
} from 'vue';
import { getUserSettings } from '@/lib/userSettings';
import type { backend } from '@wailsjs/models';

import { GetDefinitionsForWord } from '@wailsjs/backend/Dictionaries';
import { useCardManager, formatDefinition } from '@/stores/CardManager';

const UserSettings = getUserSettings();
const cardManager = useCardManager();

const definitions = ref<backend.DictionaryDefinition[]>([]);
const selectedDefs = ref<backend.DictionaryDefinition[]>([]);
let allowUpdate = false;

const props = defineProps<{
  chinese?: boolean
  english?: boolean
}>();

function updateDefinition(definitions: backend.DictionaryDefinition[]) {
  cardManager.updateDefinition(definitions,
    props.english ? 'english' : 'chinese');
}

async function calculateDefault() {
  if (definitions.value.length === 1) {
    selectedDefs.value.push(...definitions.value);
    updateDefinition(definitions.value);
  }
}

watch(selectedDefs, async () => {
  if (!allowUpdate) {
    return;
  }
  updateDefinition(selectedDefs.value);
  if (props.english
    ? UserSettings.CardCreation.AutoAdvanceEnglish
    : UserSettings.CardCreation.AutoAdvanceChinese
  ) {
    cardManager.nextStep();
  }
});

onBeforeMount(async () => {
  if (cardManager.originalValues) {
    const originalDefinitions =
      props.english
        ? cardManager.originalValues.englishDefn
        : cardManager.originalValues.chineseDefn;
    if (originalDefinitions) {
      selectedDefs.value.push(...originalDefinitions);
      definitions.value.push(...originalDefinitions);
    }
  }

  definitions.value.push(...(await GetDefinitionsForWord(
    cardManager.word,
    props.english ? 'english' : 'chinese',
  )).filter((def) => {
    return !definitions.value.some(
      (other) => {
        return def.definition === other.definition &&
        def.pronunciation === other.pronunciation;
      },
    );
  }));

  if (!cardManager.originalValues &&
      (props.english
        ? UserSettings.CardCreation.PopulateEnglish
        : UserSettings.CardCreation.PopulateChinese)
  ) {
    calculateDefault();
  }
  allowUpdate = true;
});

</script>
