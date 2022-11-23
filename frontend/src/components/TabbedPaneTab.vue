<template>
  <div v-if="isActive" :title="title">
    <slot />
  </div>
</template>

<script setup>
import { inject, watch, ref } from 'vue';

const activeTab = inject('activeTab');
const isActive = ref(true);

const props = defineProps({
  title: {
    type: String,
    default: 'TabbedPaneTab',
  },
});

watch(activeTab, () => {
  isActive.value = activeTab.value === props.title;
});
</script>
