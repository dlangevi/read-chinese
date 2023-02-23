<template>
  <!-- The button to open modal -->
  <label for="card-settings" class="cursor-pointer">
    <settings-sharp class="h-6 w-6" />
  </label>

  <!-- Put this part before </body> tag -->
  <input id="card-settings" type="checkbox" class="modal-toggle">
  <label for="card-settings" class="modal cursor-pointer">
    <div class="modal-box relative w-1/2 max-w-5xl">
      <p class="text-xl">
        Card Creation Settings
      </p>
      <div class="divider" />
      <label
        for="card-settings"
        class="btn-sm btn-circle btn absolute right-2 top-2"
      >✕</label>
      <div class="grid grid-cols-2 gap-8">
        <settings-checkbox
          v-for="(value, key) in UserSettings.CardCreation"
          :key="key"
          :setting="ComponentTable[key]"
          :initial-value="value"
        />
      </div>
    </div>
  </label>
</template>

<script lang="ts" setup>
import { SettingsSharp } from '@vicons/ionicons5';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import { ComponentTable, getUserSettings } from '@/lib/userSettings';

const UserSettings = getUserSettings();

</script>
