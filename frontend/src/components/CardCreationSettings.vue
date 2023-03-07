<template>
  <div>
    <settings-sharp
      class="h-6 w-6 cursor-pointer"
      @click="() => cardSettingsModal = true"
    />
    <div
      :class="['modal', {'modal-open': cardSettingsModal}]"
      @click="() => cardSettingsModal = false"
    >
      <div
        class="modal-box relative w-1/2 max-w-5xl"
        @click.stop
      >
        <close-circle-sharp
          class="absolute right-2 top-2 h-6 w-6 cursor-pointer"
          @click="() => cardSettingsModal = false"
        />
        <h3 class="text-xl">
          Card Creation Settings
        </h3>
        <div class="divider" />
        <div class="grid grid-cols-2 gap-8">
          <settings-checkbox
            v-for="(value, key) in ComponentTable.CardCreation"
            :key="key"
            :setting="value"
            :initial-value="getValue(key)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { SettingsSharp, CloseCircleSharp } from '@vicons/ionicons5';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import { ComponentTable, getUserSettings } from '@/lib/userSettings';

const UserSettings = getUserSettings();
const CardCreation = UserSettings.CardCreation;

function getValue(key :string) {
  return CardCreation[key as keyof typeof CardCreation] as boolean;
}
const cardSettingsModal = ref(false);
</script>
