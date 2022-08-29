import { defineStore } from 'pinia';

export const useCardQueue = defineStore('CardQueue', {
  state: () => ({ wordList: [] }),
  getters: {
    words: (state) => state.wordList,
  },
  actions: {
    addWord(word) {
      this.wordList.push(word);
    },
  },
});
