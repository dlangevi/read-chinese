<template>
  <div class="col-span-2 flex items-center gap-4">
    <span>
      {{ setting.label }}
    </span>
    <select
      v-model="currentValue"
      class="select-primary select"
      @change="submitChange"
    >
      <option
        v-for="option in options"
        :key="option"
        :value="option"
      >
        {{ option }}
      </option>
    </select>
    <div
      v-if="setting.tooltip"
      class="tooltip tooltip-right"
      :data-tip="setting.tooltip"
    >
      <information-circle class="h-6 w-6" />
    </div>
  </div>
</template>

<script lang="ts" setup>
// TODO have some kind of hover popup with more info
import { ref } from 'vue';
import { InformationCircle } from '@vicons/ionicons5';

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
    type: String,
    required: true,
  },
});

const currentValue = ref(props.initialValue);
const options = ref<string[]>([]);

async function fetchFields() {
  if (props.setting.dataSource !== undefined) {
    props.setting.dataSource().then((loaded : string[]) => {
      if (loaded !== undefined) {
        options.value = loaded;
      } else {
        console.log('Failed to fetch data for selector');
        setTimeout(fetchFields, 1000);
      }
    }).catch((err : string) => {
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
