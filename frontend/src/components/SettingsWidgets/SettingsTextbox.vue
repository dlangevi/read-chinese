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
    <button
      v-if="readonly"
      class="btn-primary btn-sm btn"
      @click="startEdit"
    >
      edit
    </button>
    <button
      v-if="!readonly"
      class="btn-primary btn-sm btn"
      @click="submitChange"
    >
      submit
    </button>

    <div
      class="tooltip tooltip-right"
      :data-tip="setting.tooltip"
    >
      <information-circle class="h-6 w-6" />
    </div>
  </n-space>
</template>

<script lang="ts" setup>
// TODO have some kind of hover popup with more info
import { ref } from 'vue';
import {
  NInput, NSpace,
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
