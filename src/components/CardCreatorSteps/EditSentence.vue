<template>
  <div class="text-3xl m-4">Pick a sentence</div>
  <div v-if="loaded && sentences.length == 0">
    No sentences found, please skip for now
  </div>
  <n-radio-group v-model:value="sentence" name="sentences">
    <n-space vertical :size="40">
    <p v-if="multiSection" class="text-4xl">From Current Book</p>
      <n-radio
        class="text-3xl"
        v-for="(sentence, i) in sentences"
        :key="i"
        :value="sentence"
        :label="sentence"
      />
    <p v-if="multiSection" class="text-4xl">From All Books</p>
      <n-radio
        class="text-3xl"
        v-for="(sentence, i) in allSentences"
        :key="i"
        :value="sentence"
        :label="sentence"
      />
    </n-space>
  </n-radio-group>
</template>

<script setup>
import {
  watch, onBeforeMount, ref, inject,
} from 'vue';
import {
  NSpace, NRadioGroup, NRadio,
} from 'naive-ui';

const UserSettings = inject('userSettings');

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
  },
});

const multiSection = props.preferBook !== undefined;

console.log('preferbook', props.preferBook);
onBeforeMount(async () => {
  if (props.preferBook) {
    sentences.value = await window.ipc.getSentencesForWord(
      props.word,
      [props.preferBook],
    );
    // TODO filter out repeats
    allSentences.value = await window.ipc.getSentencesForWord(
      props.word,
    );
  } else {
    sentences.value = await window.ipc.getSentencesForWord(
      props.word,
    );
  }
  loaded.value = true;
});

</script>
