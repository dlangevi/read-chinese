<template>
  <div class="text-3xl m-4">Pick a definition</div>
  <n-checkbox-group v-model:value="definition" name="definitions">
    <n-space vertical>
      <n-checkbox
        class="text-3xl"
        v-for="(definition, i) in definitions"
        :key="i"
        :value="definition.definition"
      >
        <span
          v-html="'[' + definition.pronunciation + '] '
            + definition.definition"
        />
      </n-checkbox>
    </n-space>
  </n-checkbox-group>
</template>

<script lang="ts" setup>
import {
  watch, onBeforeMount, ref,
} from 'vue';
import {
  NSpace, NCheckboxGroup, NCheckbox,
} from 'naive-ui';
import { getUserSettings } from '@/UserSettings';
import type { DictionaryType, DictionaryEntry } from '@/shared/types';

const UserSettings = getUserSettings();

const emit = defineEmits(['updateDefinition']);
const definitions = ref<DictionaryEntry[]>([]);
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
  emit('updateDefinition', selectedDefinitions, autoAdvance);
});

const props = defineProps<{
  word: string
  type: DictionaryType
}>();

onBeforeMount(async () => {
  definitions.value = await window.nodeIpc.getDefinitionsForWord(
    props.word,
    props.type,
  );
});

</script>
