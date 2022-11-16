<template>
  <div class="col-span-2 flex gap-2">
    <span>
      {{ setting.label }}
    </span>
    <n-slider
      v-model:value="value"
      :step="1"
      @update:value="submitChange"
    />
    <n-input-number
      v-model:value="value"
      size="small"
      @update:value="submitChange"
    />
    <div
      class="tooltip tooltip-right"
      :data-tip="setting.tooltip"
    >
      <information-circle class="h-6 w-6" />
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  NSlider, NInputNumber,
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
