<template>
  <div class="text-3xl m-4">Pick a definition</div>
  <n-checkbox-group v-model:value="definition" name="definitions">
    <n-space vertical>
      <n-checkbox
        class="text-3xl"
        v-for="(definition, i) in definitions"
        :key="i"
        :value="definition"
      >
        <span
          v-html="'[' + definition.pronunciation + '] '
            + definition.definition"
        />
      </n-checkbox>
    </n-space>
  </n-checkbox-group>
</template>

<script setup>
import {
  watch, onBeforeMount, ref, inject,
} from 'vue';
import {
  NSpace, NCheckboxGroup, NCheckbox,
} from 'naive-ui';
import { UserSettingsKey } from '../../../shared/types';

const UserSettings = inject(UserSettingsKey);

const emit = defineEmits(['updateDefinition']);
const definitions = ref([]);
const definition = ref(null);

watch(definition, async () => {
  // TODO either rename this option or have a select based on type
  const autoAdvance = await (
    UserSettings.CardCreation.AutoAdvanceEnglish.read()
  );
  console.log('update definition');
  emit('updateDefinition', definition.value, autoAdvance);
});

const props = defineProps({
  word: {
    type: String,
    required: true,
  },
  type: {
    type: String,
    required: true,
  },
});

onBeforeMount(async () => {
  definitions.value = await window.ipc.getDefinitionsForWord(
    props.word,
    props.type,
  );
});

</script>
