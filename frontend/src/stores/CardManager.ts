import { defineStore } from 'pinia';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import {
  backend,
} from '@wailsjs/models';

export const StepState = {
  // The start state for any field
  EMPTY: 'empty',
  // Currently not used
  AUTOFILL: 'autofill',
  // When entries have been selected, we move them to the card preview
  // This should not trigger the 'ready' state, so 'preview' gets upgrade
  // To 'filled' only after nextStep has been called
  PREVIEW: 'preview',
  // If nextStep is called on an empty field,
  // we assume they dont want to fill it
  SKIPPED: 'skipped',
  FILLED: 'filled',
} as const;

export type StepState = typeof StepState[keyof typeof StepState]

type StateMap = {
  [key in StepsEnum]: StepState
};

export const useCardManager = defineStore('CardManager', {
  state: () => {
    return {
      steps: [] as StepsEnum[],
      currentStep: 'sentence' as StepsEnum,
      currentStepIndex: 0,
      // Flow describes when the user is progressing as normal, either
      // auto advancing or clicking next step. In this state, when we
      // advance to a field that has been auto filled, we skip over it
      flow: true,
      sentenceSource: '',
      stepsState: {} as StateMap,
      originalValues: new backend.Fields(),
      newValues: new backend.Fields(),
      note: new backend.RawAnkiNote(),
    };
  },
  getters: {

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
    ready: (state) => {
      return state.flow &&
      Object.keys(state.stepsState).length > 0 &&
      Object.values(state.stepsState).every(state => {
        return (
          state !== StepState.EMPTY &&
          state !== StepState.PREVIEW
        );
      });
    },

  },
  actions: {
    loadCard(ankiCard : backend.RawAnkiNote,
      includeChinese : boolean,
      hasImageApi:boolean) {
      // Resets the ui (Does it?)
      this.flow = false;
      this.currentStep = StepsEnum.NONE;
      this.steps = [];

      this.steps = [
        StepsEnum.SENTENCE,
        StepsEnum.ENGLISH,
        // toggle this based on user settings
        ...(includeChinese ? [StepsEnum.CHINESE] : []),
        ...(hasImageApi ? [StepsEnum.IMAGE] : []),
      ];
      this.steps.forEach((step) => {
        this.stepsState[step] = StepState.EMPTY;
      });

      this.note = ankiCard;
      this.originalValues = backend.Fields.createFrom(this.note.fields);
      this.newValues = backend.Fields.createFrom();
      this.newValues.word = this.originalValues.word;
      this.sentenceSource = '';

      this.currentStep = StepsEnum.SENTENCE;
      this.currentStepIndex = 0;
      this.flow = true;
    },

    updateSentence(sentence: backend.Sentence) {
      this.newValues.sentence = sentence.sentence;
      console.log('setting source', sentence.source);
      this.sentenceSource = sentence.source;
      this.stepsState[StepsEnum.SENTENCE] = StepState.PREVIEW;
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
        this.stepsState[StepsEnum.ENGLISH] = StepState.PREVIEW;
      } else {
        this.newValues.chineseDefn = definitions;
        this.stepsState[StepsEnum.CHINESE] = StepState.PREVIEW;
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
      this.stepsState[StepsEnum.IMAGE] = StepState.PREVIEW;
    },

    changeStep(step: StepsEnum) {
      this.flow = false;
      this.currentStep = step;
      this.currentStepIndex = this.steps.indexOf(step);
    },

    stepState(step: StepsEnum) {
      return this.stepsState[step];
    },

    previousStep() {
      this.flow = false;
      if (this.currentStepIndex === 0) {
        return;
      }

      this.currentStepIndex -= 1;
      this.currentStep = this.steps[this.currentStepIndex];
    },
    nextStep() {
      const currentState = this.stepsState[this.currentStep];
      if (currentState === StepState.EMPTY) {
        this.stepsState[this.currentStep] = StepState.SKIPPED;
      } else if (currentState === StepState.PREVIEW) {
        this.stepsState[this.currentStep] = StepState.FILLED;
      }
      console.log('current state', this.stepsState);
      if (this.currentStepIndex + 1 === this.steps.length) {
        // We were on the last step
        return;
      }
      this.currentStepIndex += 1;
      this.currentStep = this.steps[this.currentStepIndex];
      if (this.flow &&
          this.stepsState[this.currentStep] !== StepState.EMPTY) {
        this.nextStep();
      }
    },

  },

});
