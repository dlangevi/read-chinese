import { defineStore } from 'pinia';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import {
  backend,
} from '@wailsjs/models';

export const useCardManager = defineStore('CardManager', {
  state: () => {
    return {
      steps: [] as StepsEnum[],
      originalValues: new backend.Fields(),
      newValues: new backend.Fields(),
      // Ideally this is a computed value
      note: new backend.RawAnkiNote(),
    };
  },
  getters: {

    currentStep: (_state) => {
      // complex logic based on various things
      return StepsEnum.SENTENCE;
    },

    // steps: (state) => {
    //   return state.steps;
    // },

    word: (state) => {
      return state.newValues.word || state.originalValues.word;
    },
    sentence: (state) => {
      return state.newValues.sentence || state.originalValues.sentence;
    },
    englishDefn: (state) => {
      return state.newValues.englishDefn || state.originalValues.englishDefn;
    },
    chineseDefn: (state) => {
      return state.newValues.chineseDefn || state.originalValues.chineseDefn;
    },
    pinyin: (state) => {
      return state.newValues.pinyin || state.originalValues.pinyin;
    },
    imageUrls: (state) => {
      return state.newValues.imageUrls || state.originalValues.imageUrls;
    },
    image64: (state) => {
      return state.newValues.image64 || state.originalValues.image64;
    },

  },
  actions: {
    loadCard(ankiCard : backend.RawAnkiNote) {
      this.steps = [
        StepsEnum.SENTENCE,
        StepsEnum.ENGLISH,
        // TODO toggle this based on user settings
        StepsEnum.CHINESE,
        StepsEnum.IMAGE,
      ];

      this.note = ankiCard;
      this.originalValues = backend.Fields.createFrom(this.note.fields);
      this.newValues = backend.Fields.createFrom();
      this.newValues.word = this.originalValues.word;
    },

    updateSentence(sentence: string) {
      this.newValues.sentence = sentence;
    },

    updateDefinition(
      newDefinitions: backend.DictionaryDefinition[],
      defType: string,
    ) {
      const definitions = newDefinitions.map(
        (def) => `[${def.pronunciation}] ${def.definition}`,
      ).join('<br>');
      if (defType === 'english') {
        this.newValues.englishDefn = definitions;
      } else {
        this.newValues.chineseDefn = definitions;
      }
      let pinyin = new Set();
      if (this.newValues.pinyin !== undefined) {
        pinyin = new Set(this.newValues.pinyin.split(', '));
      }
      newDefinitions.forEach((def) => {
        const pronunciation = def.pronunciation.replace(/\s/g, '');
        pinyin.add(pronunciation);
      });
      pinyin.delete('');
      this.newValues.pinyin = [...pinyin].join(', ');
    },

    updateImages(newImages: backend.ImageInfo[]) {
      this.newValues.imageUrls = newImages.map((image) => image.thumbnailUrl);
    },

  },

});
