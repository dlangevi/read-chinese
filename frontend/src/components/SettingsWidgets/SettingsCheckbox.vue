<template>
  <div class="flex items-center gap-4">
    <label
      class="flex items-center gap-4"
      :for="setting.value"
    >
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
      {{ setting.label }}
    </label>
    <div
      v-if="setting.tooltip"
      class="tooltip tooltip-top"
      :data-tip="setting.tooltip"
    >
      <information-circle class="h-6 w-6" />
    </div>
  </div>
</template>

<script lang="ts" setup>
// TODO have some kind of hover popup with more info
import { InformationCircle } from '@vicons/ionicons5';
import { ref } from 'vue';

const emit = defineEmits(['update']);

const props = defineProps({
  setting: {
    type: Object,
    required: false,
    default: () => ({
      label: 'Error: passed undefined setting',
      value: true,
    }),
  },
  initialValue: {
    type: Boolean,
    required: true,
  },
});

const isChecked = ref(props.initialValue);

function submitChange(event: Event) {
  const checked = (event.target as HTMLInputElement).checked;
  console.log('is now', checked);
  props.setting.write(checked);
  emit('update', checked);
}

</script>
