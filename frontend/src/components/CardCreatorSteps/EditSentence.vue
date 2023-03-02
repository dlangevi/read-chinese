<template>
  <div>
    <div class="m-4 text-3xl">
      Pick a sentence
    </div>
    <div v-if="loaded && (allSentences.length + sentences.length) == 0">
      No sentences found, please skip for now
    </div>
    <div v-if="props.preferBook">
      <p class="text-4xl">
        From Current Book
      </p>
      <div
        v-for="(sen, i) in sentences"
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
    <div>
      <p class="text-4xl">
        From All Books
      </p>

      <div
        v-for="(sen, i) in allSentences"
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

const sentences = ref<backend.Sentence[]>([]);
const allSentences = ref<backend.Sentence[]>([]);
const sentence = ref<backend.Sentence>(backend.Sentence.createFrom());
const loaded = ref(false);

watch(sentence, async () => {
  cardManager.updateSentence(sentence.value);
  const autoAdvance = UserSettings.CardCreation.AutoAdvanceSentence;
  if (autoAdvance) {
    cardManager.nextStep();
  }
});

onBeforeMount(async () => {
  if (props.preferBook) {
    sentences.value = await GetSentencesForWord(
      cardManager.word,
      [props.preferBook],
    );
  }
  // filter out repeats (TODO do this be passing a negative filter to
  // GetSentencesForWord ?)
  const unfilteredAll = await GetSentencesForWord(cardManager.word, []);
  allSentences.value = unfilteredAll.filter(
    (sen) => sentences.value.some(
      (other) => sen.sentence === other.sentence));
  loaded.value = true;
});

</script>
