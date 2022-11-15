<template>
  <n-card>
    <h2 class="text-xl font-bold">
      Hanzi: {{ ankiCard.fields.word }}
    </h2>
    <div>
      <n-button
        type="info"
        size="tiny"
        class="inline m-2"
        @click="$emit('change-step', StepsEnum.SENTENCE)"
      >
        Edit
      </n-button>
      <h2 class="inline text-l font-bold">
        Sentence:
      </h2>
      <span> {{ ankiCard.fields.sentence }}</span>
    </div>
    <div>
      <n-button
        type="info"
        size="tiny"
        class="m-2"
        @click="$emit('change-step', StepsEnum.ENGLISH)"
      >
        Edit
      </n-button>
      <h2 class="inline text-l font-bold">
        Definition:
      </h2>
      <p>
        <!-- Definitions can have html formatting in them -->
        <span v-html="ankiCard.fields.englishDefn" />
      </p>
    </div>
    <div v-if="enableChinese">
      <n-button
        type="info"
        size="tiny"
        class="m-2"
        @click="$emit('change-step', StepsEnum.CHINESE)"
      >
        Edit
      </n-button>
      <h2 class="inline text-l font-bold">
        ChineseDefinition:
      </h2>
      <p>
        <!-- Definitions can have html formatting in them -->
        <span v-html="ankiCard.fields.chineseDefn" />
      </p>
    </div>
    <div>
      <n-button
        type="info"
        size="tiny"
        class="m-2"
        @click="$emit('change-step', StepsEnum.IMAGE)"
      >
        Edit
      </n-button>
      <h2 class="inline text-l font-bold">
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
import { NCard, NButton } from 'naive-ui';
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
