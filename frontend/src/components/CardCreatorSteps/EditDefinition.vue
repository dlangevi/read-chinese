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
          v-model="definition"
          class="checkbox"
          :value="def.definition"
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
import { storeToRefs } from 'pinia';
import { getUserSettings } from '@/lib/userSettings';
import type { backend } from '@wailsjs/models';

import { GetDefinitionsForWord } from '@wailsjs/backend/Dictionaries';
import { useCardManager } from '@/stores/CardManager';

const UserSettings = getUserSettings();

const cardManager = useCardManager();
const { word } = storeToRefs(cardManager);

const definitions = ref<backend.DictionaryDefinition[]>([]);
const definition = ref<string[]>([]);

watch(definition, async () => {
  const selected = new Set<string>(definition.value.values());
  const selectedDefinitions = definitions.value.filter(
    (def) => selected.has(def.definition),
  );
  updateDefinition(selectedDefinitions);
  const autoAdvance = await (
    UserSettings.CardCreation.AutoAdvanceEnglish.read()
  );
  if (autoAdvance) {
    cardManager.nextStep();
  }
});

watch(word, () => {
  loadData();
});

const props = defineProps<{
  type: string
}>();

function updateDefinition(definitions: backend.DictionaryDefinition[]) {
  cardManager.updateDefinition(definitions, props.type);
}

async function calculateDefault() {
  const definitions = await GetDefinitionsForWord(
    cardManager.word,
    props.type,
  );
  if (definitions.length === 1) {
    updateDefinition(definitions);
  }
}

async function loadData() {
  let autoFill : boolean;
  if (props.type === 'english') {
    autoFill = UserSettings.CardCreation.PopulateEnglish.read();
  } else {
    autoFill = UserSettings.CardCreation.PopulateChinese.read();
  }
  definitions.value = await GetDefinitionsForWord(
    cardManager.word,
    props.type,
  );

  if (autoFill) {
    calculateDefault();
  }
}
onBeforeMount(loadData);

</script>
