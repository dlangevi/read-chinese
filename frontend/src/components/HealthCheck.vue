<template>
  <div class="">
    <div v-if="allComplete">
      <router-link class="btn" to="/BookLibrary">Lets get started</router-link>
    </div>
    <div v-else>
      <h2 class="text-2xl font-extrabold text-white">
        Before Getting Started
      </h2>
      <ul class="list-outside list-disc space-y-2">
        <li
          v-for="check in checks"
          :key="check.description"
          class="items-center space-x-2"
        >
          <span :class="{ 'line-through': check.checkResult }">
            {{ check.description }}
          </span>
          <div
            v-if="!check.checkResult"
            class="btn-primary btn-sm btn"
            @click="check.buttonAction"
          >
            {{ check.buttonText }}
          </div>
        </li>
      </ul>
    </div>
    <div
      :class="['modal', {'modal-open': ankiInfo}]"
      @click="() => ankiInfo = false"
    >
      <div class="modal-box">
        <h3 class="text-xl font-extrabold">
          How to connect Anki
        </h3>
        <p>
          For this application to work you need to have Anki opened along with
          the Anki-Connect plugin installed.
        </p>
        <ul class="mt-4 list-inside list-disc text-lg">
          <li>
            <a
              href="#"
              class="text-primary hover:text-primary-focus"
              @click="BrowserOpenURL('https://apps.ankiweb.net/')"
              @click.stop
            >
              Download anki from here
            </a>
          </li>
          <li class="mt-2">
            <a
              href="#"
              class="text-primary hover:text-primary-focus"
              @click="BrowserOpenURL('https://foosoft.net/projects/anki-connect/')"
              @click.stop
            >
              Setup up the Anki-Connect plugin
            </a>
          </li>
        </ul>
      </div>
    </div>
    <div
      :class="['modal', {'modal-open': dictInfo}]"
      @click="() => dictInfo = false"
    >
      <div class="modal-box" @click.stop>
        <dictionaries-list />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, reactive, computed, ref } from 'vue';
import { BrowserOpenURL } from '@runtime/runtime';
import { HealthCheck as bookHealth } from '@wailsjs/backend/bookLibrary';
import { HealthCheck as dictHealth } from '@wailsjs/backend/Dictionaries';
import { HealthCheck as ankiHealth } from '@wailsjs/backend/ankiInterface';
import { ImportCalibreBooks } from '@wailsjs/backend/Calibre';
import DictionariesList
  from '@/components/SettingsWidgets/DictionariesList.vue';
const ankiInfo = ref(false);
const dictInfo = ref(false);

const Check = {
  BOOKLIBRARY: 'books',
  DICTIONARY: 'dictionary',
  ANKIAVALIABLE: 'ankiavaliable',
  ANKICONFIGURED: 'ankiconfigured',
} as const;

const checks = reactive({
  [Check.DICTIONARY]: {
    description: 'Have at least one Dictionary installed',
    checkAction: dictHealth,
    checkResult: false,
    buttonText: 'Add dictionary',
    buttonAction: () => { dictInfo.value = true; },
  },
  [Check.BOOKLIBRARY]: {
    description: 'Import at least one book',
    checkAction: bookHealth,
    checkResult: false,
    buttonText: 'Sync Calibre',
    buttonAction: async () => {
      const err = await ImportCalibreBooks();
      console.log(err);
    },
  },
  [Check.ANKIAVALIABLE]: {
    description: 'Anki is avaliable through anki-connect',
    checkAction: ankiHealth,
    checkResult: false,
    buttonText: 'Get Help',
    buttonAction: () => { ankiInfo.value = true; },
  },
  [Check.ANKICONFIGURED]: {
    description: 'Configure the names of anki fields',
    // Currently only I use this so its true because of
    // hardcoded anki settings
    checkAction: async () => true,
    checkResult: false,
    buttonText: 'Configure Anki',
    buttonAction: () => true,
  },
});

const allComplete = computed(() => {
  return Object.values(checks).every(check => check.checkResult);
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
