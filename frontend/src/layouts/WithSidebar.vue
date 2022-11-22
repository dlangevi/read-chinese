<template>
  <div class="drawer drawer-end">
    <input id="my-drawer" type="checkbox" class="drawer-toggle">
    <div class="drawer-content">
      <slot />
    </div>
    <div class="drawer-side">
      <label for="my-drawer" class="drawer-overlay" />
      <ul class="menu w-80 gap-4 bg-base-100 p-4 text-base-content">
        <slot name="sidebar" />
      </ul>
    </div>
  </div>
</template>
