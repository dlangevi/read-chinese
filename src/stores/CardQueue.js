import { defineStore } from 'pinia';

export const useCardQueue = defineStore('CardQueue', {
  state: () => ({ wordList: [] }),
  getters: {
    words: (state) => state.wordList,
  },
  actions: {
    async addWord(word) {
      this.wordList.push(word);
    },
    clearWords() {
      this.wordList.splice(0);
    },
  },
});
