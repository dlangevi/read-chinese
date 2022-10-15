<template>
  <div class="text-3xl m-4">Pick a sentence</div>
  <div v-if="loaded && (allSentences.length + sentences.length) == 0">
    No sentences found, please skip for now
  </div>
  <n-radio-group v-model:value="sentence" name="sentences">
    <div v-if="singleBook">
      <n-space vertical :size="40">
        <p class="text-4xl">From Current Book</p>
        <n-radio
          class="text-3xl"
          v-for="(sentence, i) in sentences"
          :key="i"
          :value="sentence"
          :label="sentence"
        />
      </n-space>
    </div>
    <div>
      <n-space vertical :size="40">
        <p class="text-4xl">From All Books</p>
        <n-radio
          class="text-3xl"
          v-for="(sentence, i) in allSentences"
          :key="i"
          :value="sentence"
          :label="sentence"
        />
      </n-space>
    </div>
  </n-radio-group>
</template>

<script setup>
import {
  watch, onBeforeMount, ref,
} from 'vue';
import {
  NSpace, NRadioGroup, NRadio,
} from 'naive-ui';
import { getUserSettings } from '@/renderer/UserSettings';

const UserSettings = getUserSettings();

const emit = defineEmits(['updateSentence']);
const sentences = ref([]);
const allSentences = ref([]);
const sentence = ref(null);
const loaded = ref(false);

watch(sentence, async () => {
  const autoAdvance = await (
    UserSettings.CardCreation.AutoAdvanceSentence.read()
  );
  emit('updateSentence', sentence.value, autoAdvance);
});

const props = defineProps({
  word: {
    type: String,
    required: true,
  },
  preferBook: {
    type: String,
    default: undefined,
  },
});

const singleBook = !!props.preferBook;
onBeforeMount(async () => {
  if (props.preferBook) {
    sentences.value = await window.ipc.getSentencesForWord(
      props.word,
      { bookIds: [props.preferBook] },
    );
  }
  // TODO filter out repeats
  allSentences.value = await window.ipc.getSentencesForWord(
    props.word,
  );
  loaded.value = true;
});

</script>
