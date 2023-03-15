<template>
  <div
    v-if="UserSettings.AzureConfig.AzureApiKey"
    class="col-span-3 row-span-2 flex flex-col gap-2"
  >
    <div
      :class="['modal', {'modal-open': voicePickerModal}]"
      @click="() => voicePickerModal = false"
    >
      <div
        class="modal-box relative w-1/2 max-w-5xl"
        @click.stop
      >
        <div class="flex">
          <h2 class="text-xl">Select New Voice</h2>
          <a
            href="#"
            class="ml-auto text-info hover:text-info-content"
            @click="
              BrowserOpenURL('https://speech.microsoft.com/portal/voicegallery')"
          >
            Visit Azure Voice Gallery for an overview of the voices
          </a>
        </div>
        <div class="grid grow grid-cols-2 gap-4">
          <div
            v-for="section, title in pickers"
            :key="title"
          >
            <div
              v-if="section.visible"
              class="flex w-full items-center p-4"
            >
              <span class="">{{ title }}</span>
              <select
                v-model="section.value"
                class="select-primary select ml-auto"
              >
                <option
                  v-for="option in section.options"
                  :key="option"
                  :value="option"
                >
                  {{ option }}
                </option>
              </select>
            </div>
          </div>
          <div class="flex items-center gap-4">
            <span>
              Speaking Speed
            </span>
            <input
              v-model.number="activeSpeed"
              class="range range-primary"
              type="range"
              min="-100"
              max="200"
              step="1"
            >
            <div>
              {{ displaySpeed(activeSpeed) }}
            </div>
          </div>
          <div class="flex items-center gap-4">
            <span>
              Pitch
            </span>
            <input
              v-model.number="activePitch"
              class="range range-primary"
              type="range"
              min="-50"
              max="50"
              step="1"
            >
            <div> {{ displayPitch(activePitch) }} </div>
          </div>
          <button class="btn-secondary btn-sm btn col-span-2" @click="addVoice">
            Add Voice
          </button>
        </div>
      </div>
    </div>
    <div class="flex items-center gap-4">
      <div class="grow text-2xl">
        Voices
      </div>
      <button
        class="btn-secondary btn w-1/4"
        @click="voicePickerModal = true"
      >
        Add New Voice
      </button>
    </div>
    <table class="table flex-col gap-4">
      <thead>
        <tr>
          <th>Voice</th>
          <th>Style</th>
          <th>Role</th>
          <th>Speed</th>
          <th>Pitch</th>
          <th />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="voice in UserSettings.AzureConfig.VoiceList"
          :key="voice.Voice"
        >
          <td> {{ voice.Voice }} </td>
          <td> {{ voice.SpeakingStyle }} </td>
          <td> {{ voice.RolePlay }} </td>
          <td> {{ displaySpeed(voice.Speed) }} </td>
          <td> {{ displayPitch(voice.Pitch) }} </td>
          <td>
            <button
              class="btn-error btn-sm btn"
              @click="removeVoice(voice)"
            >
              Remove Voice
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref, onBeforeMount, computed, watch } from 'vue';
import { BrowserOpenURL } from '@runtime/runtime';
import {
  GetVoices,
} from '@wailsjs/backend/TextToSpeech';
import {
  AddVoice,
  RemoveVoice,
} from '@wailsjs/backend/UserSettings';
import {
  backend,
} from '@wailsjs/models';

import { getUserSettings } from '@/lib/userSettings';
const UserSettings = getUserSettings();

const voices = ref<backend.TTSVoice[]>([]);

const activeSpeed = ref(0);
const activePitch = ref(0);

function displaySpeed(speed:number) {
  return ((speed + 100) / 100).toFixed(2);
}

function displayPitch(pitch:number) {
  return ((pitch + 50) / 50).toFixed(2);
}

const activeVoiceString = ref('');
const voicePickerModal = ref(false);

computed(() => {
  const matching = voices.value.filter(
    voice => voice.DisplayName === activeVoiceString.value);

  if (matching.length > 0) {
    return matching[0];
  }
  return backend.TTSVoice.createFrom();
});

const pickers = reactive({
  Locale: {
    options: [] as string[],
    value: '',
    visible: true,
  },
  SpeakingStyle: {
    options: [] as string[],
    value: 'General',
    visible: false,
  },
  Voice: {
    options: [] as string[],
    value: '',
    visible: true,
  },
  RolePlay: {
    options: [] as string[],
    value: 'Default',
    visible: false,
  },
},
);

watch(() => pickers.Locale.value, () => {
  pickers.Voice.options = voices.value
    .filter((voice) => voice.Locale === pickers.Locale.value)
    .map((voice) => voice.DisplayName);

  pickers.Voice.value = pickers.Voice.options[0];
});

watch(() => pickers.Voice.value, () => {
  const activeVoice = voices.value.find(
    voice => voice.DisplayName === pickers.Voice.value);
  if (!activeVoice) {
    return;
  }

  pickers.SpeakingStyle.options = [
    'General',
    ...(activeVoice.StyleList
      ? activeVoice.StyleList
      : []),
  ];
  pickers.SpeakingStyle.value = pickers.SpeakingStyle.options[0];
  pickers.SpeakingStyle.visible =
    activeVoice.StyleList !== undefined;

  pickers.RolePlay.options = [
    'Default',
    ...(activeVoice.RolePlayList
      ? activeVoice.RolePlayList
      : []),
  ];

  pickers.RolePlay.value = pickers.RolePlay.options[0];
  pickers.RolePlay.visible =
    activeVoice.RolePlayList !== undefined;
});

async function addVoice() {
  const activeVoice = voices.value.find(
    voice => voice.DisplayName === pickers.Voice.value);
  if (!activeVoice) {
    return;
  }

  const voice = new backend.Voice();
  voice.Locale = pickers.Locale.value;
  voice.Voice = activeVoice.ShortName;
  voice.SpeakingStyle = pickers.SpeakingStyle.value;
  if (voice.SpeakingStyle === 'General') {
    voice.SpeakingStyle = '';
  }
  voice.RolePlay = pickers.RolePlay.value;
  if (voice.RolePlay === 'Default') {
    voice.RolePlay = '';
  }
  voice.Speed = activeSpeed.value;
  voice.Pitch = activePitch.value;
  await AddVoice(voice);
  voicePickerModal.value = false;
}

function removeVoice(voice: backend.Voice) {
  RemoveVoice(voice);
}

async function loadVoices() {
  if (!UserSettings.AzureConfig.AzureApiKey) {
    return;
  }
  voices.value = await GetVoices();
  pickers.Locale.options = [
    ...new Set(voices.value.map((voice) => voice.Locale)),
  ];
  pickers.Locale.value = pickers.Locale.options[0];
}

watch(() => UserSettings.AzureConfig.AzureApiKey, () => {
  loadVoices();
});

onBeforeMount(async () => {
  loadVoices();
});

</script>
