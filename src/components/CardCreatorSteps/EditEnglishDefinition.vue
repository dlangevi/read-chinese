<template>
  <div class="text-3xl m-4">Pick a definition</div>
  <n-checkbox-group v-model:value="definition" name="definitions">
    <n-space vertical>
      <n-checkbox
        class="text-3xl"
        v-for="(definition, i) in definitions"
        :key="i"
        :value="definition"
        :label="
        '[' + definition.pronunciation + '] ' +
        definition.definition"
      />
    </n-space>
  </n-checkbox-group>
</template>

<script setup>
import { watch, onBeforeMount, ref } from 'vue';
import {
  NSpace, NCheckboxGroup, NCheckbox,
} from 'naive-ui';
import UserSettings from '@/userSettings';

const emit = defineEmits(['updateDefinition']);
const definitions = ref([]);
const definition = ref(null);

watch(definition, async () => {
  const autoAdvance = await UserSettings.AutoAdvanceEnglish.read();
  emit('updateDefinition', definition.value, autoAdvance);
});

const props = defineProps({
  word: {
    type: String,
    required: true,
  },
});

onBeforeMount(async () => {
  definitions.value = await window.ipc.getDefinitionsForWord(props.word);
});

</script>
