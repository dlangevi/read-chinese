<template>
  <n-space class="col-span-2">
    <span>
      {{ setting.label }}
    </span>
    <n-input
      :type="type"
      show-password-on="mousedown"
      class="min-w-8"
      :placeholder="setting.label"
      :default-value="initialValue"
      :maxlength="32"
      :readonly="readonly"
      @update:value="updateValue"
    />
    <n-button
      v-if="readonly"
      @click="startEdit"
    >
      edit
    </n-button>
    <n-button
      v-if="!readonly"
      @click="submitChange"
    >
      submit
    </n-button>
    <n-tooltip
      v-if="setting.tooltip"
      placement="right"
      trigger="hover"
    >
      <template #trigger>
        <n-icon size="20">
          <information-circle />
        </n-icon>
      </template>
      <span>{{ setting.tooltip }}</span>
    </n-tooltip>
  </n-space>
</template>

<script lang="ts" setup>
// TODO have some kind of hover popup with more info
import { ref } from 'vue';
import {
  NTooltip, NIcon, NInput, NButton, NSpace,
} from 'naive-ui';
import { InformationCircle } from '@vicons/ionicons5';

const type = ref<'password' | 'text'>('password');
const readonly = ref(true);
const props = defineProps({
  setting: {
    type: Object,
    required: true,
  },
});

const initialValue = await props.setting.read();
const currentValue = ref(initialValue);

function updateValue(newValue: string) {
  currentValue.value = newValue;
}
function startEdit() {
  readonly.value = false;
  type.value = 'text';
}

function submitChange() {
  // TODO actually save
  console.log(currentValue.value);
  props.setting.write(currentValue.value);
  readonly.value = true;
  type.value = 'password';
}

</script>
