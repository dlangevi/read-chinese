<template>
  <div>
    <div class="text-3xl m-4">
      Pick a sentence
    </div>
    <div v-if="loaded && (allSentences.length + sentences.length) == 0">
      No sentences found, please skip for now
    </div>
    <n-radio-group
      v-model:value="sentence"
      name="sentences"
    >
      <div v-if="singleBook">
        <n-space
          vertical
          :size="40"
        >
          <p class="text-4xl">
            From Current Book
          </p>
          <n-radio
            v-for="(sen, i) in sentences"
            :key="i"
            class="text-3xl"
            :value="sen"
            :label="sen"
          />
        </n-space>
      </div>
      <div>
        <n-space
          vertical
          :size="40"
        >
          <p class="text-4xl">
            From All Books
          </p>
          <n-radio
            v-for="(sen, i) in allSentences"
            :key="i"
            class="text-3xl"
            :value="sen"
            :label="sen"
          />
        </n-space>
      </div>
    </n-radio-group>
  </div>
</template>

<script lang="ts" setup>
import {
  watch, onBeforeMount, ref,
} from 'vue';
import {
  NSpace, NRadioGroup, NRadio,
} from 'naive-ui';
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
  // TODO filter out repeats
  allSentences.value = await GetSentencesForWord(props.word, []);
  loaded.value = true;
});

</script>
