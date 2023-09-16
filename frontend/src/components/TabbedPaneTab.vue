<template>
  <div v-show="isActive" :tab-title="tabTitle" class="h-full">
    <slot />
  </div>
</template>

<script setup>
import { inject, watch, ref } from 'vue';

const activeTab = inject('activeTab');
const isActive = ref(true);

const props = defineProps({
  tabTitle: {
    type: String,
    required: true,
  },
});

watch(activeTab, () => {
  isActive.value = activeTab.value === props.tabTitle;
});
</script>
