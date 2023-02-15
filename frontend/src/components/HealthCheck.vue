<template>
  <div class="overflow-x-auto">
    <table class="table w-full">
      <thead>
        <tr>
          <th>Completed</th>
          <th>Description</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>
            {{ checks[Check.DICTIONARY].checkResult }}
          </td>
          <td>Has at least one Dictionary installed</td>
          <td><div class="btn">Do it</div></td>
        </tr>
        <tr>
          <td>
            {{ checks[Check.BOOKLIBRARY].checkResult }}
          </td>
          <td>Has at least one book imported</td>
          <td><div class="btn">Do it</div></td>
        </tr>
        <tr>
          <td>
            {{ checks[Check.ANKIAVALIABLE].checkResult }}
          </td>
          <td>Anki is avaliable through anki-connect</td>
          <td><div class="btn">Do it</div></td>
        </tr>
        <tr>
          <td>
            {{ checks[Check.ANKICONFIGURED].checkResult }}
          </td>
          <td>The correct fields have been configured</td>
          <td><div class="btn">Do it</div></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, reactive } from 'vue';
import { HealthCheck as bookHealth } from '@wailsjs/backend/bookLibrary';
import { HealthCheck as dictHealth } from '@wailsjs/backend/Dictionaries';
import { HealthCheck as ankiHealth } from '@wailsjs/backend/ankiInterface';

const Check = {
  BOOKLIBRARY: 'books',
  DICTIONARY: 'dictionary',
  ANKIAVALIABLE: 'ankiavaliable',
  ANKICONFIGURED: 'ankiconfigured',
} as const;

const checks = reactive({
  [Check.DICTIONARY]: {
    description: 'Add a dictionary',
    checkAction: dictHealth,
    checkResult: false,
  },
  [Check.BOOKLIBRARY]: {
    description: 'Add a book',
    checkAction: bookHealth,
    checkResult: false,
  },
  [Check.ANKIAVALIABLE]: {
    description: 'Connect to Anki',
    checkAction: ankiHealth,
    checkResult: false,
  },
  [Check.ANKICONFIGURED]: {
    description: 'Configure Anki cards',
    // Currently only I use this so its true because of
    // hardcoded anki settings
    checkAction: async () => true,
    checkResult: false,
  },
});

function recheck() {
  Object.values(checks).forEach(async (check) => {
    check.checkResult = await check.checkAction();
  });
}
recheck();
const checkInterval = setInterval(recheck, 1000);
onBeforeUnmount(() => {
  clearInterval(checkInterval);
});

</script>
