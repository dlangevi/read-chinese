<template>
  <n-card>
    <h2 class="text-xl font-bold">
      Hanzi: {{ ankiCard.fields.word }}
    </h2>
    <div>
      <button
        class="btn btn-primary btn-xs m-2 inline"
        @click="$emit('change-step', StepsEnum.SENTENCE)"
      >
        Edit
      </button>
      <h2 class="text-l inline font-bold">
        Sentence:
      </h2>
      <span> {{ ankiCard.fields.sentence }}</span>
    </div>
    <div>
      <button
        class="btn btn-primary btn-xs m-2"
        @click="$emit('change-step', StepsEnum.ENGLISH)"
      >
        Edit
      </button>
      <h2 class="text-l inline font-bold">
        Definition:
      </h2>
      <p>
        <!-- Definitions can have html formatting in them -->
        <span v-html="ankiCard.fields.englishDefn" />
      </p>
    </div>
    <div v-if="enableChinese">
      <button
        class="btn btn-primary btn-xs m-2"
        @click="$emit('change-step', StepsEnum.CHINESE)"
      >
        Edit
      </button>
      <h2 class="text-l inline font-bold">
        ChineseDefinition:
      </h2>
      <p>
        <!-- Definitions can have html formatting in them -->
        <span v-html="ankiCard.fields.chineseDefn" />
      </p>
    </div>
    <div>
      <button
        class="btn btn-primary btn-xs m-2"
        @click="$emit('change-step', StepsEnum.IMAGE)"
      >
        Edit
      </button>
      <h2 class="text-l inline font-bold">
        Images:
      </h2>
      <div class="flex gap-1">
        <img
          v-for="image in ankiCard.fields.imageUrls"
          :key="image"
          class="max-h-24 w-auto"
          :src="image"
          alt="image for word"
        >
      </div>
    </div>
  </n-card>
</template>

<script lang="ts" setup>
import { NCard } from 'naive-ui';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import { getUserSettings } from '@/lib/userSettings';

const UserSettings = getUserSettings();

const enableChinese = UserSettings.Dictionaries.EnableChinese.read();

defineEmits(['change-step']);
defineProps({
  ankiCard: {
    type: Object,
    required: true,
  },
});

</script>
