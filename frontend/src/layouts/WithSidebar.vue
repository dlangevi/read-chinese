<template>
  <div class="drawer-mobile drawer">
    <input id="my-drawer" type="checkbox" class="drawer-toggle">
    <div class="drawer-content">
      <slot />
    </div>
    <aside class="drawer-side">
      <label for="my-drawer" class="drawer-overlay" />
      <ul class="menu w-80 gap-4 bg-base-300 p-4 text-base-content">
        <slot name="sidebar" />
      </ul>
    </aside>
  </div>
</template>
