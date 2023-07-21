<template>
  <div
    v-if="finishedIntital"
    :class="[
      'rounded-lg',
      {'bg-base-100': !allComplete },
      'p-4'
    ]"
  >
    <div v-if="stillChecking">
      Loading ...
    </div>
    <div v-else-if="allComplete">
      <router-link class="btn" to="/BookLibrary">Lets get started</router-link>
    </div>
    <div v-else>
      <div class="flex gap-4">
        <h2 class="text-2xl font-extrabold text-base-content">
          Before Getting Started
        </h2>
        <button
          class="
          btn-primary
          btn"
          @click="firstFailure?.buttonAction"
        >
          {{ firstFailure?.buttonText }}
        </button>
      </div>
      <ul class="list-inside list-disc space-y-2">
        <li
          v-for="check in checks"
          :key="check.description"
          class="items-center space-x-2"
        >
          <span :class="{ 'line-through': check.passes}">
            {{ check.description }}
            <span
              v-if="!check.passes"
              class="font-bold text-error"
            >({{ check.checkResult }})</span>
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
          How to install and connect to Anki
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
              How to setup up the Anki-Connect plugin
            </a>
          </li>
        </ul>
      </div>
    </div>
    <div
      :class="['modal', {'modal-open': dictInfo}]"
      @click="() => dictInfo = false"
    >
      <div
        class="modal-box"
        @click.stop
      >
        <dictionaries-list />
      </div>
    </div>
    <div
      :class="['modal', {'modal-open': ankiConfigure}]"
      @click="() => ankiConfigure = false"
    >
      <div
        v-if="checks.AnkiAvaliable.passes"
        class="modal-box flex w-4/5 max-w-full flex-col gap-4"
        @click.stop
      >
        <settings-selector
          :setting="ComponentTable.AnkiConfig.ActiveDeck"
          :initial-value="UserSettings.AnkiConfig.ActiveDeck"
        />
        <settings-selector
          :setting="ComponentTable.AnkiConfig.ActiveModel"
          :initial-value="UserSettings.AnkiConfig.ActiveModel"
        />

        <model-manager />
      </div>
    </div>
    <div
      :class="['modal', {'modal-open': bookImport}]"
      @click="() => bookImport = false"
    >
      <div
        class="modal-box flex w-4/5 max-w-full flex-col gap-4"
        @click.stop
      >
        <book-importer />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, watch, reactive, computed, ref } from 'vue';
import { BrowserOpenURL } from '@runtime/runtime';
import { ComponentTable, getUserSettings } from '@/lib/userSettings';

import { HealthCheck as bookHealth } from '@wailsjs/backend/bookLibrary';
import { HealthCheck as dictHealth } from '@wailsjs/backend/Dictionaries';
import {
  HealthCheck as ankiHealth,
  ConfigurationCheck as ankiConfigured,
} from '@wailsjs/backend/ankiInterface';

import SettingsSelector from
  '@/components/SettingsWidgets/SettingsSelector.vue';
import BookImporter from
  '@/components/BookImporter.vue';

import DictionariesList
  from '@/components/SettingsWidgets/DictionariesList.vue';
import ModelManager
  from '@/components/SettingsWidgets/ModelManager.vue';

const UserSettings = getUserSettings();
const ankiInfo = ref(false);
const dictInfo = ref(false);
const bookImport = ref(false);
const ankiConfigure = ref(false);

type HealthCheckInfo = {
  buttonText: string,
  description: string
  checkAction: () => Promise<void>,
  buttonAction: () => void,
  passes?: boolean,
  checkResult?: string,
}

const checks = reactive<{
  Dictionary: HealthCheckInfo,
  BookLibrary: HealthCheckInfo,
  AnkiAvaliable: HealthCheckInfo,
  AnkiConfigured: HealthCheckInfo,
}>({
  Dictionary: {
    buttonText: 'Add dictionary',
    description: 'Have at least one Dictionary installed',
    checkAction: dictHealth,
    buttonAction: () => { dictInfo.value = true; },
  },
  BookLibrary: {
    buttonText: 'Sync Calibre',
    description: 'Import at least one book',
    checkAction: bookHealth,
    buttonAction: async () => { bookImport.value = true; },
  },
  AnkiAvaliable: {
    buttonText: 'How to setup Anki',
    description: 'Anki is avaliable through anki-connect',
    checkAction: ankiHealth,
    buttonAction: () => { ankiInfo.value = true; },
  },
  AnkiConfigured: {
    buttonText: 'Configure Anki',
    description: 'Configure the names of anki fields',
    checkAction: async () => {
      if (checks.AnkiAvaliable.passes) {
        return ankiConfigured();
      } else {
        throw checks.AnkiAvaliable.checkResult;
      }
    },
    buttonAction: () => { ankiConfigure.value = true; },
  },
});

type checkKeys = keyof typeof checks

watch(
  () => [
    checks.Dictionary.passes,
    checks.AnkiAvaliable.passes,
    checks.AnkiConfigured.passes,
    checks.BookLibrary.passes,
  ],
  () => {
    if (checks.Dictionary.passes) {
      dictInfo.value = false;
    }
    if (checks.BookLibrary.passes) {
      bookImport.value = false;
    }
    if (checks.AnkiAvaliable.passes) {
      ankiInfo.value = false;
    }
    if (checks.AnkiConfigured.passes) {
      ankiConfigure.value = false;
    }
  });

watch(
  () =>
    checks.AnkiAvaliable.passes,
  () => {
    if (checks.AnkiAvaliable.passes) {
      recheck(['AnkiConfigured']);
    }
  });

const allComplete = computed(() => {
  return Object.values(checks).every(
    check => check.passes);
});

const firstFailure = computed(() => {
  return Object.values(checks).find(
    check => check.passes !== true);
});

const stillChecking = computed(() => {
  return Object.values(checks).some(
    check => check.passes === undefined);
});

function recheck(listOfChecks : checkKeys[] = [
  'Dictionary',
  'BookLibrary',
  'AnkiAvaliable',
  'AnkiConfigured']) {
  return Promise.all((listOfChecks).map(async (checkKey) => {
    const check = checks[checkKey];
    return check.checkAction()
      .then(() => { check.passes = true; })
      .catch(errMsg => {
        check.passes = false;
        check.checkResult = errMsg;
      });
  }));
}

const finishedIntital = ref(false);
recheck().finally(async () => {
  finishedIntital.value = true;
});

const checkInterval = setInterval(recheck, 1000);

onBeforeUnmount(() => {
  clearInterval(checkInterval);
});

</script>
