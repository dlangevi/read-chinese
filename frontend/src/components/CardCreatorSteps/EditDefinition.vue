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
        :for="def.definition"
      >
        <input
          :id="def.definition"
          v-model="selectedDefs"
          class="checkbox"
          :value="def"
          type="checkbox"
          name="definitions"
        >
        <span v-html="'[' + def.pronunciation + '] ' + def.definition" />
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
import { useCardManager } from '@/stores/CardManager';

const UserSettings = getUserSettings();
const cardManager = useCardManager();

const definitions = ref<backend.DictionaryDefinition[]>([]);
const selectedDefs = ref<backend.DictionaryDefinition[]>([]);

const props = defineProps<{
  chinese?: boolean
  english?: boolean
}>();

watch(selectedDefs, async () => {
  updateDefinition(selectedDefs.value);
  if (props.english
    ? UserSettings.CardCreation.AutoAdvanceEnglish
    : UserSettings.CardCreation.AutoAdvanceChinese
  ) {
    cardManager.nextStep();
  }
});

function updateDefinition(definitions: backend.DictionaryDefinition[]) {
  cardManager.updateDefinition(definitions,
    props.english ? 'english' : 'chinese');
}

async function calculateDefault() {
  if (definitions.value.length === 1) {
    updateDefinition(definitions.value);
  }
}

onBeforeMount(async () => {
  definitions.value = await GetDefinitionsForWord(
    cardManager.word,
    props.english ? 'english' : 'chinese',
  );

  if (props.english
    ? UserSettings.CardCreation.PopulateEnglish
    : UserSettings.CardCreation.PopulateChinese
  ) {
    calculateDefault();
  }
});

</script>
