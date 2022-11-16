<template>
  <div>
    <input
      :id="setting.value"
      v-model="isChecked"
      type="checkbox"
      :label="setting.label"
      :disabled="setting.disabled"
      @input="submitChange"
    >
    <label
      v-if="setting.tooltip"
      :for="setting.value"
    >
      {{ setting.label }}
    </label>
    <div
      class="tooltip"
      :data-tip="setting.tooltip"
    >
      <information-circle class="h-4 w-4" />
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

const isChecked = await props.setting.read();

function submitChange(event: Event) {
  const checked = (event.target as HTMLInputElement).checked;
  console.log('is now', checked);
  props.setting.write(checked);
  emit('update', checked);
}

</script>
