<template>
  <n-button
    text
    type="info"
    @click="openSettings"
  >
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
      <div class="grid grid-cols-2 gap-8">
        <settings-checkbox
          v-for="content in UserSettings.CardCreation"
          :key="content.name"
          :setting="content"
        />
      </div>
    </Suspense>
  </n-modal>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import {
  NIcon, NButton, NModal,
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
