<template>
  <n-button
    @click="openSettings"
    text
    type="info">
    <template #icon>
      <n-icon size="20">
        <settings-sharp />
      </n-icon>
    </template>
  </n-button>
  <n-modal
    v-model:show="showSettings"
    :closable="true"
    class="w-1/2"
    preset="card"
  >
    <template #header>
      <p class="text-xl">
        Card Creation Settings
      </p>
    </template>

    <Suspense>
      <n-grid :y-gap="8" :cols="2">
        <n-gi
          v-for="content in UserSettings.CardCreation"
          :key="content"
        >
          <settings-checkbox :setting="content" />
        </n-gi>
      </n-grid>
    </Suspense>

  </n-modal>
</template>

<script setup>
import { ref } from 'vue';
import {
  NIcon, NButton, NModal, NGrid, NGi,
} from 'naive-ui';
import { SettingsSharp } from '@vicons/ionicons5';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import { getUserSettings } from '@/lib/userSettings';

const UserSettings = getUserSettings();

const showSettings = ref(false);

const openSettings = () => {
  showSettings.value = true;
};

</script>
