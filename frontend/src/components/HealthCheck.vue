<template>
  <div class="">
    <div v-if="allComplete">
      <router-link class="btn" to="/BookLibrary">Lets get started</router-link>
    </div>
    <div v-else>
      <div class="flex gap-4">
        <h2 class="text-2xl font-extrabold text-white">
          Before Getting Started
        </h2>
        <div
          class="btn-primary btn"
          @click="nextAction"
        >
          {{ nextActionText }}
        </div>
      </div>
      <ul class="list-outside list-disc space-y-2">
        <li
          v-for="check in checks"
          :key="check.description"
          class="items-center space-x-2"
        >
          <span :class="{ 'line-through': check.checkResult == ''}">
            {{ check.description }}
          </span>
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
    <div
      :class="['modal', {'modal-open': ankiConfigure}]"
      @click="() => ankiConfigure = false"
    >
      <div
        v-if="checks.ANKIAVALIABLE.checkResult === ''"
        class="modal-box flex w-4/5 max-w-full flex-col gap-4"
        @click.stop
      >
        <settings-selector
          :setting="ComponentTable.ActiveDeck"
          :initial-value="UserSettings.AnkiConfig.ActiveDeck"
        />
        <settings-selector
          :setting="ComponentTable.ActiveModel"
          :initial-value="UserSettings.AnkiConfig.ActiveModel"
        />

        <model-manager />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, reactive, computed, ref } from 'vue';
import { BrowserOpenURL } from '@runtime/runtime';
import { HealthCheck as bookHealth } from '@wailsjs/backend/bookLibrary';
import { HealthCheck as dictHealth } from '@wailsjs/backend/Dictionaries';
import {
  HealthCheck as ankiHealth,
  ConfigurationCheck as ankiConfigured,
} from '@wailsjs/backend/ankiInterface';
import { ImportCalibreBooks } from '@wailsjs/backend/Calibre';
import SettingsSelector from
  '@/components/SettingsWidgets/SettingsSelector.vue';
import DictionariesList
  from '@/components/SettingsWidgets/DictionariesList.vue';
import ModelManager
  from '@/components/SettingsWidgets/ModelManager.vue';
import { ComponentTable, getUserSettings } from '@/lib/userSettings';
import { useLoader } from '@/lib/loading';
const loader = useLoader();
const UserSettings = getUserSettings();
const ankiInfo = ref(false);
const dictInfo = ref(false);
const ankiConfigure = ref(false);

const checks = reactive({
  DICTIONARY: {
    description: 'Have at least one Dictionary installed',
    checkAction: dictHealth,
    checkResult: 'not checked yet',
    buttonText: 'Add dictionary',
    buttonAction: () => { dictInfo.value = true; },
  },
  BOOKLIBRARY: {
    description: 'Import at least one book',
    checkAction: bookHealth,
    checkResult: 'not checked yet',
    buttonText: 'Sync Calibre',
    buttonAction: async () => {
      await loader.withLoader(ImportCalibreBooks, 'Importing calibre');
    },
  },
  ANKIAVALIABLE: {
    description: 'Anki is avaliable through anki-connect',
    checkAction: ankiHealth,
    checkResult: 'not checked yet',
    buttonText: 'Get Help',
    buttonAction: () => { ankiInfo.value = true; },
  },
  ANKICONFIGURED: {
    description: 'Configure the names of anki fields',
    // Currently only I use this so its true because of
    // hardcoded anki settings
    checkAction: async () => {
      const health = await ankiHealth();
      if (health !== '') {
        return health;
      } else {
        return ankiConfigured();
      }
    },
    checkResult: 'not checked yet',
    buttonText: 'Configure Anki',
    buttonAction: () => { ankiConfigure.value = true; },
  },
});

const allComplete = computed(() => {
  return Object.values(checks).every(check => check.checkResult === '');
});

const nextActionText = computed(() => {
  if (checks.DICTIONARY.checkResult !== '') {
    return checks.DICTIONARY.buttonText;
  }
  if (checks.BOOKLIBRARY.checkResult !== '') {
    return checks.BOOKLIBRARY.buttonText;
  }
  if (checks.ANKIAVALIABLE.checkResult !== '') {
    return checks.ANKIAVALIABLE.buttonText;
  }
  if (checks.ANKICONFIGURED.checkResult !== '') {
    return checks.ANKICONFIGURED.buttonText;
  }
  return 'error';
});
const nextAction = computed(() => {
  if (checks.DICTIONARY.checkResult !== '') {
    return checks.DICTIONARY.buttonAction;
  }
  if (checks.BOOKLIBRARY.checkResult !== '') {
    return checks.BOOKLIBRARY.buttonAction;
  }
  if (checks.ANKIAVALIABLE.checkResult !== '') {
    return checks.ANKIAVALIABLE.buttonAction;
  }
  if (checks.ANKICONFIGURED.checkResult !== '') {
    return checks.ANKICONFIGURED.buttonAction;
  }
  return () => {};
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
