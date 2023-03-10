<template>
  <div class="h-full">
    <div class="tabs tabs-boxed mb-4">
      <a
        v-for="(tab, i) in tabs"
        :key="tab.getAttribute('tab-title') as string"
        :class="['tab', {'tab-active': i == selected}]"
        @click="selectTab(i)"
      >
        {{ tab.getAttribute('tab-title') }}
      </a>
    </div>
    <div ref="children" class="h-full">
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, provide } from 'vue';
const selected = ref(0);
const tabs = ref<HTMLElement[]>([]);
const children = ref<HTMLDivElement | null>(null);
const activeTab = ref('');

provide('activeTab', activeTab);

onMounted(() => {
  if (children.value === null) {
    console.error('No provided tabs:', children);
    return;
  }
  tabs.value = Array.from(children.value.children) as HTMLElement[];
  selectTab(0);
});

function selectTab(i:number) {
  selected.value = i;
  activeTab.value = tabs.value[i].getAttribute('tab-title') as string;
}
</script>
