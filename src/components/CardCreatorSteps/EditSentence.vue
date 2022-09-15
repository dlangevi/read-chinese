<template>
  <div class="text-3xl m-4">Pick a sentence</div>
  <div v-if="loaded && sentences.length == 0">
    No sentences found, please skip for now
  </div>
  <n-radio-group v-model:value="sentence" name="sentences">
    <n-space vertical :size="40">
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
import UserSettings from '@/userSettings';

const emit = defineEmits(['updateSentence']);
const sentences = ref([]);
const sentence = ref(null);
const loaded = ref(false);

watch(sentence, async () => {
  const autoAdvance = await UserSettings.AutoAdvanceSentence.read();
  emit('updateSentence', sentence.value, autoAdvance);
});

const props = defineProps({
  word: {
    type: String,
    required: true,
  },
});

onBeforeMount(async () => {
  sentences.value = await window.ipc.getSentencesForWord(props.word);
  loaded.value = true;
});

</script>
