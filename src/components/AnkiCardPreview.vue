<template>
  <n-card>
    <h2 class="text-xl font-bold">Hanzi: {{ankiCard.fields.Hanzi.value}}</h2>
    <p>
      <span class="text-l font-bold">Sentence:</span>
      {{ankiCard.fields.ExampleSentence.value}}
      <n-button type=info size=tiny class="inline"
        @click="$emit('changeStep', StepsEnum.SENTENCE)">Edit</n-button>
    </p>
    <h2 class="text-l font-bold">Definition:</h2>
    <p>
      <span v-html="ankiCard.fields.EnglishDefinition.value"/>
      <n-button type=info size=tiny
        @click="$emit('changeStep', StepsEnum.ENGLISH)">Edit</n-button>
    </p>
    <h2 class="text-l font-bold">ChineseDefinition:</h2>
    <p>
      <span v-html="ankiCard.fields.ChineseDefinition.value"/>
      <n-button type=info size=tiny
        @click="$emit('changeStep', StepsEnum.CHINESE)">Edit</n-button>
    </p>
  </n-card>
</template>

<script setup>
import { NCard, NButton } from 'naive-ui';
import StepsEnum from '@/components/CardCreatorSteps/StepsEnum';

const props = defineProps({
  ankiCard: {
    type: Object,
    required: true,
  },
});

console.log(props.ankiCard);

</script>
