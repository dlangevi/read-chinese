<template>
  <img @click="markKnown(this.params.data)" class="
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

<script>
export default {
  data() {
    return {
    };
  },
  methods: {
    markKnown: (rowData) => {
      console.log(rowData);
      // TODO
      // ipc.addNewWord(...)
      // remove this row from data somehow
    },

  },
  beforeMount() {
    // console.log(this.params.data);
  },
};
</script>
