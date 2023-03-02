<template>
  <div class="flex items-center gap-4">
    <button
      :class="['btn-primary btn-sm btn',
               {'invisible': cardManager.currentStepIndex == 0}]"
      @click="cardManager.previousStep()"
    >
      Previous Step
    </button>
    <ul class="steps">
      <li
        v-for="(step, i) in cardManager.steps"
        :key="step"
        :data-content="i"
        :class="[
          'step',
          getClass(step),
          {'after:ring after:ring-4 after:ring-accent':
            cardManager.currentStep === step},
        ]"
      >
        {{ step }}
      </li>
    </ul>
    <button
      v-if="cardManager.steps.length > 0"
      class="btn-primary btn-sm btn"
      :class="['btn-primary btn-sm btn']"
      @click="cardManager.nextStep()"
    >
      Next Step
    </button>
  </div>
</template>

<script setup>
import { useCardManager, StepState } from '@/stores/CardManager';
const cardManager = useCardManager();

function getClass(step) {
  const state = cardManager.stepsState[step];
  const mapping = {
    [StepState.EMPTY]: 'step-neutral',
    [StepState.PREVIEW]: 'step-accent',
    [StepState.FILLED]: 'step-success',
    [StepState.SKIPPED]: 'step-success',
  };
  return mapping[state];
}
</script>
