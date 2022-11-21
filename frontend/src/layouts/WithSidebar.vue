<template>
  <div class="flex h-full w-full flex-wrap">
    <div
      id="sidebar"
      class="flex shrink grow basis-1/6 flex-col gap-4 p-4"
    >
      <slot name="sidebar" />
    </div>
    <div
      id="content"
      class="shrink grow basis-5/6"
    >
      <slot />
    </div>
  </div>
</template>
