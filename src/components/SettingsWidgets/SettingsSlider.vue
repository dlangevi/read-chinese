<template>
  <div class="col-span-2 flex gap-2">
    <span>
      {{setting.label}}
    </span>
    <n-slider
      @update:value="submitChange"
      v-model:value="value"
      :step="1" />
    <n-input-number
      @update:value="submitChange"
      v-model:value="value"
      size="small" />
    <n-tooltip
      v-if="setting.tooltip"
      placement="right"
      trigger="hover">
      <template #trigger>
        <n-icon size="20">
          <information-circle />
        </n-icon>
      </template>
      <span> {{setting.tooltip}}</span>
    </n-tooltip>
  </div>
</template>

<script setup lang="ts">
import {
  NTooltip, NIcon, NSlider, NInputNumber,
} from 'naive-ui';
import { ref } from 'vue';
import { InformationCircle } from '@vicons/ionicons5';

const emit = defineEmits(['update']);

const props = defineProps({
  setting: {
    type: Object,
    required: true,
  },
});

const value = ref(props.setting.read());

function submitChange(newValue:number | null) {
  props.setting.write(newValue);
  emit('update', newValue);
}

</script>
