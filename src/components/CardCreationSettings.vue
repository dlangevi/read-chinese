<template>
  <n-button
    @click="openSettings"
    text type="info">
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
        <n-gi>
          <settings-checkbox :setting="UserSettings.PopulateEnglish" />
        </n-gi>
        <n-gi>
          <settings-checkbox :setting="UserSettings.AutoAdvanceSentence" />
        </n-gi>
        <n-gi>
          <settings-checkbox :setting="UserSettings.AutoAdvanceEnglish" />
        </n-gi>
        <n-gi>
          <settings-checkbox :setting="UserSettings.AutoAdvanceCard" />
        </n-gi>
        <n-gi>
          <settings-checkbox :setting="UserSettings.EnableChinese" />
        </n-gi>
        <n-gi>
          <settings-checkbox :setting="UserSettings.PopulateChinese" />
        </n-gi>
        <n-gi>
          <settings-checkbox :setting="UserSettings.GenerateTermAudio" />
        </n-gi>
        <n-gi>
          <settings-checkbox :setting="UserSettings.GenerateSentenceAudio" />
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
import SettingsCheckbox from '@/components/SettingsCheckbox.vue';
import UserSettings from '@/userSettings';

const showSettings = ref(false);

const openSettings = () => {
  showSettings.value = true;
};

</script>
