import { defineStore } from 'pinia';

export const useCardQueue = defineStore('CardQueue', {
  state: () => ({ wordList: [] }),
  getters: {
    words: (state) => state.wordList,
  },
  actions: {
    // Callback runs on word submition
    // TODO? also have a callback for word failure?
    async addWord(word, action, callback) {
      this.wordList.push({
        word,
        action,
        callback,
      });
    },
    clearFront() {
      this.wordList.shift();
    },
    clearWords() {
      this.wordList.splice(0);
    },
  },
});

export const ActionsEnum = Object.freeze({
  CREATE: 'create',
  MODIFY: 'modify',
});
