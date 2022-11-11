<template>
  <div>
    <n-checkbox
      :value="setting.value"
      :label="setting.label"
      @update:checked="submitChange"
      :default-checked="isChecked"
      :disabled="setting.disabled"
    />
    <n-tooltip v-if="setting.tooltip" placement="right" trigger="hover">
      <template #trigger>
        <n-icon size="20">
          <information-circle />
        </n-icon>
      </template>
      <span> {{setting.tooltip}}</span>
    </n-tooltip>
  </div>
</template>

<script lang="ts" setup>
// TODO have some kind of hover popup with more info
import {
  NCheckbox, NTooltip, NIcon,
} from 'naive-ui';
import { InformationCircle } from '@vicons/ionicons5';

const emit = defineEmits(['update']);

const props = defineProps({
  setting: {
    type: Object,
    required: true,
  },
});

const isChecked = await props.setting.read();

function submitChange(checked: boolean) {
  props.setting.write(checked);
  emit('update', checked);
}

</script>
