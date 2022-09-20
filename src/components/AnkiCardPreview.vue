<template>
  <n-card>
    <h2 class="text-xl font-bold">Hanzi: {{ankiCard.fields.word}}</h2>
    <div>
      <n-button type=info size=tiny class="inline m-2"
        @click="$emit('changeStep', StepsEnum.SENTENCE)">Edit</n-button>
      <h2 class="inline text-l font-bold">Sentence:</h2>
      <span> {{ankiCard.fields.sentence}}</span>
    </div>
    <div>
      <n-button type=info size=tiny class="m-2"
        @click="$emit('changeStep', StepsEnum.ENGLISH)">Edit</n-button>
      <h2 class="inline text-l font-bold">Definition:</h2>
      <p>
        <span v-html="ankiCard.fields.englishDefn"/>
      </p>
    </div>
    <div v-if="enableChinese">
      <n-button type=info size=tiny class="m-2"
        @click="$emit('changeStep', StepsEnum.CHINESE)">Edit</n-button>
      <h2 class="text-l font-bold">ChineseDefinition:</h2>
      <p >
        <span v-html="ankiCard.fields.chinesDefn"/>
      </p>
    </div>
  </n-card>
</template>

<script setup>
import { NCard, NButton } from 'naive-ui';
import StepsEnum from '@/components/CardCreatorSteps/StepsEnum';
import { inject } from 'vue';

const UserSettings = inject('userSettings');

const enableChinese = UserSettings.Dictionaries.EnableChinese.read();
console.log(`enablechinese ${enableChinese}`);

defineProps({
  ankiCard: {
    type: Object,
    required: true,
  },
});

</script>
