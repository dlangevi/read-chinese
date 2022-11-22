<template>
  <div>
    <div class="m-4 text-3xl">
      Pick a sentence
    </div>
    <div v-if="loaded && (allSentences.length + sentences.length) == 0">
      No sentences found, please skip for now
    </div>
    <div v-if="singleBook">
      <p class="text-4xl">
        From Current Book
      </p>
      <div
        v-for="(sen, i) in sentences"
        :key="i"
        class="text-3xl"
      >
        <label class="label cursor-pointer" :for="sen">
          <input
            :id="sen"
            v-model="sentence"
            class="radio"
            :value="sen"
            type="radio"
            name="sentences"
          >
          <span class="label-text">{{ sen }}</span>
        </label>
      </div>
    </div>
    <div>
      <p class="text-4xl">
        From All Books
      </p>

      <div
        v-for="(sen, i) in allSentences"
        :key="i"
        class="text-3xl"
      >
        <label class="label cursor-pointer" :for="sen">
          <input
            :id="sen"
            v-model="sentence"
            class="radio"
            :value="sen"
            type="radio"
            name="sentences"
          >
          <span class="label-text">{{ sen }}</span>
        </label>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  watch, onBeforeMount, ref,
} from 'vue';
import { getUserSettings } from '@/lib/userSettings';
import { GetSentencesForWord } from '@wailsjs/backend/Generator';

const UserSettings = getUserSettings();

const emit = defineEmits(['update-sentence']);
const sentences = ref<string[]>([]);
const allSentences = ref<string[]>([]);
const sentence = ref(null);
const loaded = ref(false);

watch(sentence, async () => {
  const autoAdvance = await (
    UserSettings.CardCreation.AutoAdvanceSentence.read()
  );
  emit('update-sentence', sentence.value, autoAdvance);
});

const props = defineProps<{
  word: string
  preferBook?: number
}>();

const singleBook = !!props.preferBook;
onBeforeMount(async () => {
  if (props.preferBook) {
    sentences.value = await GetSentencesForWord(
      props.word,
      [props.preferBook],
    );
  }
  // filter out repeats
  allSentences.value = (await GetSentencesForWord(props.word, []))
    .filter((sen) => {
      return sentences.value.indexOf(sen) === -1;
    });
  loaded.value = true;
});

</script>
