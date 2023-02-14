<template>
  <div class="container mx-auto px-4">
    <div class="text-center">
      <h2 class="mt-5 text-center">
        Welcome to read-chinese
      </h2>
      <p>
        Designed to help manage a dynamic flashcard
        library to aid in reading books, lets get started by making sure
        some beggining steps have been completed
      </p>
      <div class="btn" @click="recheck">Recheck</div>
    </div>
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
  </div>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
import { HealthCheck as bookHealth } from '@wailsjs/backend/bookLibrary';
import { HealthCheck as dictHealth } from '@wailsjs/backend/Dictionaries';
import { HealthCheck as ankiHealth } from '@wailsjs/backend/ankiInterface';

const Priority = {
  REQUIRED: 'required',
  OPTIONAL: 'optional',
} as const;

type PriorityType = typeof Priority[keyof typeof Priority];
type Step = {
  priority: PriorityType;
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
    priority: Priority.REQUIRED,
    description: 'Add a dictionary',
    checkAction: dictHealth,
    checkResult: false,
    workButton: true,
  },
  {
    priority: Priority.REQUIRED,
    description: 'Add a book',
    checkAction: bookHealth,
    checkResult: false,
    workButton: true,
  },
  {
    priority: Priority.REQUIRED,
    description: 'Connect to Anki',
    checkAction: ankiHealth,
    checkResult: false,
    workButton: true,
  },
  {
    priority: Priority.REQUIRED,
    description: 'Configure Anki cards',
    // Currently only I use this so its true because of
    // hardcoded anki settings
    checkAction: async () => true,
    checkResult: false,
    workButton: true,
  },
  {
    priority: Priority.OPTIONAL,
    description: 'Import word list',
    checkAction: async () => true,
    checkResult: false,
    workButton: true,
  },
  {
    priority: Priority.OPTIONAL,
    description: 'Mark words Known',
    checkAction: async () => false,
    checkResult: false,
    workButton: true,
  },
  {
    priority: Priority.OPTIONAL,
    description: 'Add azure API key',
    checkAction: async () => true,
    checkResult: false,
    workButton: true,
  },
  {
    priority: Priority.OPTIONAL,
    description: 'Add azure Image key',
    checkAction: async () => false,
    checkResult: false,
    workButton: true,
  },
]);

function recheck() {
  steps.forEach(async (step) => {
    console.log('before step', step.description);
    step.checkResult = await step.checkAction();
    console.log('after step', step.description, step.checkResult);
  });
}
recheck();

</script>
