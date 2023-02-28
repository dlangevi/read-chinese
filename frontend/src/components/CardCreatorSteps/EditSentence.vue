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

const cardManager = useCardManager();

const sentences = ref<backend.Sentence[]>([]);
const allSentences = ref<backend.Sentence[]>([]);
const sentence = ref<backend.Sentence>(backend.Sentence.createFrom());
const loaded = ref(false);
const UserSettings = getUserSettings();

watch(sentence, async () => {
  console.log('Setting sentence', sentence.value);
  cardManager.updateSentence(sentence.value);
  const autoAdvance = UserSettings.CardCreation.AutoAdvanceSentence;
  if (autoAdvance) {
    cardManager.nextStep();
  }
});

const props = defineProps<{
  preferBook?: number
}>();

const singleBook = !!props.preferBook;
async function loadData() {
  loaded.value = false;
  sentences.value = [];
  allSentences.value = [];
  if (props.preferBook) {
    sentences.value = await GetSentencesForWord(
      cardManager.word,
      [props.preferBook],
    );
  }
  // filter out repeats
  allSentences.value = (await GetSentencesForWord(cardManager.word, []))
    .filter((sen) => {
      return sentences.value.indexOf(sen) === -1;
    });
  loaded.value = true;
}
onBeforeMount(loadData);

</script>
