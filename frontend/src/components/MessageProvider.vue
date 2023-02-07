<template>
  <div>
    <div class="toast-center toast toast-top">
      <div
        v-if="api.shouldShow"
        :class="[
          'alert',
          {'alert-info': api.currentType === 'info'},
          {'alert-success': api.currentType === 'success'},
          {'alert-warning': api.currentType === 'warning'},
          {'alert-error': api.currentType === 'error'},
        ]"
      >
        <div>
          <span>{{ api.message }}</span>
        </div>
      </div>
    </div>
    <slot />
  </div>
</template>

<script lang="ts" setup>
import { useMessage } from '@/lib/messages';

const api = useMessage();

</script>
