<template>
  <n-radio-group v-model:value="sentence" name="sentences">
    <n-space vertical>
      <n-radio
        class="text-3xl"
        v-for="(sentence, i) in sentences"
        :key="i"
        :value="sentence"
        :label="sentence"
      />
    </n-space>
  </n-radio-group>
</template>

<script setup>
import { watch, onBeforeMount, ref } from 'vue';
import {
  NSpace, NRadioGroup, NRadio,
} from 'naive-ui';

const emit = defineEmits(['updateSentence']);
const sentences = ref([]);
const sentence = ref(null);

watch(sentence, () => {
  emit('updateSentence', sentence.value);
});

const props = defineProps({
  word: {
    type: String,
    required: true,
  },
});

onBeforeMount(async () => {
  console.log(props);
  sentences.value = await window.ipc.getSentencesForWord(props.word);
});

</script>
