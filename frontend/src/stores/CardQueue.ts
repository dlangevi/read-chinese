import { defineStore } from 'pinia';
import type { LoadOptions } from './CardManager';

type WordlistEntry = {
  options:LoadOptions
  preferBook?:number
  callback?: () => void
};

export const useCardQueue = defineStore('CardQueue', {
  state: () => ({
    wordList: [] as WordlistEntry[],
  }),
  getters: {
    words: (state) => state.wordList,
  },
  actions: {
    // Callback runs on word submition
    // TODO? also have a callback for word failure?
    async addWord(
      cardOptions:LoadOptions,
      callback?: () => void,
      preferBook?:number,
    ) {
      this.wordList.push({
        options: cardOptions,
        callback,
        preferBook,
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
