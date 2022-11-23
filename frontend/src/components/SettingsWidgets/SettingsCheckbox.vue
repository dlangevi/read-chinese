<template>
  <div class="flex gap-2">
    <input
      :id="setting.value"
      v-model="isChecked"
      type="checkbox"
      class="checkbox-primary checkbox"
      :label="setting.label"
      :disabled="setting.disabled"
      @input="submitChange"
    >
    <!-- TODO how to align the label vertically? -->
    <label
      v-if="setting.tooltip"
      class="flex-1"
      :for="setting.value"
    >
      {{ setting.label }}
    </label>
    <div
      class="tooltip tooltip-right"
      :data-tip="setting.tooltip"
    >
      <information-circle class="h-6 w-6" />
    </div>
  </div>
</template>

<script lang="ts" setup>
// TODO have some kind of hover popup with more info
import { InformationCircle } from '@vicons/ionicons5';

const emit = defineEmits(['update']);

const props = defineProps({
  setting: {
    type: Object,
    required: true,
  },
});

const isChecked = props.setting.read();

function submitChange(event: Event) {
  const checked = (event.target as HTMLInputElement).checked;
  console.log('is now', checked);
  props.setting.write(checked);
  emit('update', checked);
}

</script>
