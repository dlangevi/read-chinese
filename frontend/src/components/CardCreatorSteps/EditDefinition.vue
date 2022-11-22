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
      <label class="label cursor-pointer" :for="def.definition">
        <input
          :id="def.definition"
          v-model="definition"
          class="checkbox"
          :value="def.definition"
          type="checkbox"
          name="definitions"
        >
        <span
          class="label-text"
          v-html="'[' + def.pronunciation + '] ' + def.definition"
        />
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

const UserSettings = getUserSettings();

const emit = defineEmits(['update-definition']);
const definitions = ref<backend.DictionaryEntry[]>([]);
const definition = ref<string[]>([]);

watch(definition, async () => {
  // TODO either rename this option or have a select based on type
  const selected = new Set<string>(definition.value.values());
  const selectedDefinitions = definitions.value.filter(
    (def) => selected.has(def.definition),
  );
  const autoAdvance = await (
    UserSettings.CardCreation.AutoAdvanceEnglish.read()
  );
  console.log('update definition');
  emit('update-definition', selectedDefinitions, autoAdvance);
});

const props = defineProps<{
  word: string
  type: string
}>();

async function calculateDefault() {
  const definitions = await GetDefinitionsForWord(
    props.word,
    props.type,
  );
  if (definitions.length === 1) {
    emit('update-definition', definitions, false);
  }
}

let autoFill : boolean;
if (props.type === 'english') {
  autoFill = UserSettings.CardCreation.PopulateEnglish.read();
} else {
  autoFill = UserSettings.CardCreation.PopulateChinese.read();
}
if (autoFill) {
  calculateDefault();
}

onBeforeMount(async () => {
  definitions.value = await GetDefinitionsForWord(
    props.word,
    props.type,
  );
});

</script>
