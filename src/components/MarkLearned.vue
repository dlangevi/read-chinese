<template>
  <img @click="markKnown" class="
    font-black
    bg-red-600
    hover:bg-green-600
    absolute
    top-1/2
    rounded-xl
    -translate-y-1/2
    h-5
    w-5" src="../assets/circle-check.svg" alt="checkmark"/>
</template>

<script setup>
const props = defineProps({
  params: {
    type: Object,
  },
});

function markKnown() {
  const rowData = props.params.data;
  // Keep with the convention of 10000 == user has claimed
  // they super known this one
  window.ipc.addWord(rowData.word, 10000);
  props.params.api.applyTransaction({
    remove: [rowData],
  });
}

</script>
