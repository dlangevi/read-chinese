<template>
  <div class="text-3xl m-4">Pick a definition</div>
  <n-radio-group v-model:value="definition" name="definitions">
    <n-space vertical>
      <n-radio
        class="text-3xl"
        v-for="(definition, i) in definitions"
        :key="i"
        :value="definition"
        :label="
        '[' + definition.pronunciation + '] ' +
        definition.definition"
      />
    </n-space>
  </n-radio-group>
</template>

<script setup>
import { watch, onBeforeMount, ref } from 'vue';
import {
  NSpace, NRadioGroup, NRadio,
} from 'naive-ui';

const emit = defineEmits(['updateDefinition']);
const definitions = ref([]);
const definition = ref(null);

watch(definition, () => {
  console.log(definition);
  emit('updateDefinition', definition.value);
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
