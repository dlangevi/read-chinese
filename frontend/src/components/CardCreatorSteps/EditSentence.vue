<template>
  <div>
    <div class="m-4 text-3xl">
      Pick a sentence
    </div>
    <div v-if="loaded && (allSentences.length + sentences.length) == 0">
      No sentences found, please skip for now
    </div>
    <div v-for="(sectionSentences, section) in sections" :key="section">
      <div v-if="sectionSentences.value.length > 0">
        <p class="text-4xl">
          {{ section }}
        </p>
        <div
          v-for="(sen, i) in sectionSentences.value"
          :key="i"
          class="text-3xl"
        >
          <label
            class="label cursor-pointer justify-start gap-2"
            :for="sen.sentence"
          >
            <input
              :id="sen.sentence"
              v-model="sentence"
              class="radio"
              :value="sen"
              type="radio"
              name="sentences"
            >
            <span>{{ sen.sentence }}</span>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  watch, onBeforeMount, ref,
} from 'vue';
import { GetSentencesForWord } from '@wailsjs/backend/Generator';
import { useCardManager } from '@/stores/CardManager';
import { getUserSettings } from '@/lib/userSettings';
import { backend } from '@wailsjs/models';

const props = defineProps<{
  preferBook?: number
}>();

const cardManager = useCardManager();
const UserSettings = getUserSettings();

const originalSentences = ref<backend.Sentence[]>([]);
const sentences = ref<backend.Sentence[]>([]);
const allSentences = ref<backend.Sentence[]>([]);
const sentence = ref<backend.Sentence>(backend.Sentence.createFrom());
const loaded = ref(false);

const sections = {
  'Original Sentence': originalSentences,
  'From Current Book': sentences,
  'From All Books': allSentences,
};

watch(sentence, async () => {
  if (!loaded.value) {
    return;
  }
  cardManager.updateSentence(sentence.value);
  const autoAdvance = UserSettings.CardCreation.AutoAdvanceSentence;
  if (autoAdvance) {
    cardManager.nextStep();
  }
});

onBeforeMount(async () => {
  const originalSentence = cardManager.originalValues?.sentence;
  if (originalSentence) {
    sentence.value = originalSentence;
    originalSentences.value.push(originalSentence);
  }

  if (props.preferBook) {
    // PreferBook and OriginalSentence cant currently both
    // be true, but filter anyways
    sentences.value = (await GetSentencesForWord(
      cardManager.word,
      [props.preferBook],
    )).filter((sen) => {
      return !originalSentences.value.some(
        (other) => sen.sentence === other.sentence);
    });
  }
  // filter out repeats (TODO do this be passing a negative filter to
  // GetSentencesForWord ?)
  const unfilteredAll = await GetSentencesForWord(cardManager.word, []);
  allSentences.value = unfilteredAll.filter((sen) => {
    return ![...sentences.value, ...originalSentences.value].some(
      (other) => sen.sentence === other.sentence);
  });

  loaded.value = true;
});

</script>
