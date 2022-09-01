<template>
  <n-card>
    <h2 class="text-xl font-bold">Hanzi: {{ankiCard.fields.Hanzi.value}}</h2>
    <p>
      <span class="text-l font-bold">Sentence:</span>
      {{ankiCard.fields.ExampleSentence.value}}
    </p>
    <h2 class="text-l font-bold">Definition:</h2>
    <p v-html="ankiCard.fields.EnglishDefinition.value"/>
    <h2 class="text-l font-bold">ChineseDefinition:</h2>
    <p v-html="ankiCard.fields.ChineseDefinition.value"/>
  </n-card>
</template>

<script setup>
import { NCard } from 'naive-ui';

const props = defineProps({
  ankiCard: {
    type: Object,
    required: true,
  },
});

console.log(props.ankiCard);

</script>
