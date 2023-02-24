<template>
  <div class="col-span-2 flex items-center gap-4">
    <span>
      {{ setting.label }}
    </span>
    <select v-model="currentValue" class="select-primary select">
      <option
        v-for="option in options"
        :key="option"
        :value="option"
      >
        {{ option }}
      </option>
    </select>
    <button
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
  </div>
</template>

<script lang="ts" setup>
// TODO have some kind of hover popup with more info
import { ref, withDefaults } from 'vue';
import { InformationCircle } from '@vicons/ionicons5';
import type { UserSetting } from '@/lib/userSettings';

const props = withDefaults(defineProps<{
  setting?: UserSetting,
  initialValue: string,
}>(), {
  setting: () => ({
    name: 'Error',
    label: 'Error: passed undefined setting',
    value: true,
    type: false,
    write: () => {},
  }),
});

const currentValue = ref(props.initialValue);
const options = ref<string[]>([]);

async function fetchFields() {
  if (props.setting.dataSource !== undefined) {
    props.setting.dataSource().then((loaded) => {
      if (loaded !== undefined) {
        options.value = loaded;
      } else {
        console.log('Failed to fetch data for selector');
        setTimeout(fetchFields, 1000);
      }
    }).catch((err) => {
      console.log(err);
      setTimeout(fetchFields, 1000);
    });
  }
}
fetchFields();

function submitChange() {
  console.log(currentValue.value);
  props.setting.write(currentValue.value);
}

</script>
