<template>
  <div class="overflow-x-auto">
    <table class="table w-full">
      <thead>
        <tr>
          <th>Step</th>
          <th>Completed</th>
          <th>Description</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(step, i) in steps"
          :key="i"
        >
          <td>
            {{ i }}
          </td>
          <td>
            {{ step.checkResult }}
          </td>
          <td>{{ step.description }}</td>
          <td><div class="btn">Do it</div></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
type Step = {
  description: string;
  checkAction: () => Promise<boolean>;
  checkResult: boolean;
  workButton: boolean;
};

// Users need to
// Required:
// * Add at least one dictionary
// * Add at least one book
// * Establish Anki connection
// * Configure Anki Cards

// Optional:
// * Import some words list
// * Establish calibre connection
// * manually mark words as known
// * Add azure API key
// * Add azure Image key
// * Check out other settings

const steps : Step[] = reactive([
  {
    description: 'Import word list',
    checkAction: async () => true,
    checkResult: false,
    workButton: true,
  },
  {
    description: 'Mark words Known',
    checkAction: async () => false,
    checkResult: false,
    workButton: true,
  },
  {
    description: 'Add azure API key',
    checkAction: async () => true,
    checkResult: false,
    workButton: true,
  },
  {
    description: 'Add azure Image key',
    checkAction: async () => false,
    checkResult: false,
    workButton: true,
  },
]);

function recheck() {
  steps.forEach(async (step) => {
    step.checkResult = await step.checkAction();
  });
}
recheck();

</script>;
