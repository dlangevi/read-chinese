<template>
  <div>
    <button
      v-if="cardManager.steps.length > 0"
      class="btn-primary btn-sm btn"
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
          {'step-neutral': cardManager.stepsState[step] === StepState.EMPTY},
          {'step-accent': cardManager.stepsState[step] === StepState.PREVIEW},
          {'step-success': cardManager.stepsState[step] === StepState.FILLED},
          {'step-success': cardManager.stepsState[step] === StepState.SKIPPED},
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
      @click="cardManager.nextStep()"
    >
      Next Step
    </button>
  </div>
</template>

<script setup>
import { useCardManager, StepState } from '@/stores/CardManager';
const cardManager = useCardManager();
</script>
