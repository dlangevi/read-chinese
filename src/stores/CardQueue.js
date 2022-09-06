import { defineStore } from 'pinia';

export const useCardQueue = defineStore('CardQueue', {
  state: () => { return { wordList: [] }; },
  getters: {
    words: (state) => { return state.wordList; },
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
