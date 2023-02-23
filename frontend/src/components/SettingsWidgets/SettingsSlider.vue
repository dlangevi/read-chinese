<template>
  <div class="col-span-2 flex gap-2">
    <span>
      {{ setting.label }}
    </span>
    <input
      v-model="value"
      class="range range-primary"
      type="range"
      min="0"
      max="110"
      step="1"
      @input="submitChange"
    >
    <div> {{ value }} </div>
    <div
      class="tooltip tooltip-right"
      :data-tip="setting.tooltip"
    >
      <information-circle class="h-6 w-6" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { InformationCircle } from '@vicons/ionicons5';

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
    type: Number,
    required: true,
  },
});

const value = ref(props.initialValue);

function submitChange(event : Event) {
  const number = parseInt((event.target as HTMLInputElement).value);
  console.log('new number', number);
  props.setting.write(number);
  emit('update', number);
}

</script>
