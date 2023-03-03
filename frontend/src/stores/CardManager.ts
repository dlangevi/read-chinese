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

// Instead of the strings we get from Anki, we want to
// manage the data using the formats we see up front
type FrontFields = {
  word: string;
  sentence?: backend.Sentence;
  englishDefn?: backend.DictionaryDefinition[];
  chineseDefn?: backend.DictionaryDefinition[];
  pinyin?: string;
  images?: backend.ImageInfo[]
}

export function transformDefinition(definition: string)
  : backend.DictionaryDefinition[] {
  // Current definitions are stored as follows
  // [pinyin] Text for definition<br>[pinyin] Text for second definition
  // We will try to match this exactly and on failure just return
  // [] Full Text of failed definition
  const splitDefinitions = definition
    .split('<br>')
    .filter((def) => def.length > 0);

  // Capture [unicode letter + numbers] remove space
  // inbetween then capture the rest
  const regex = /\[([\p{L}\s\d]*)\]\s+(.+)/u;
  const allMatch = splitDefinitions.every((def) => regex.test(def));
  if (!allMatch) {
    return [{
      definition,
      pronunciation: undefined,
    }];
  }

  return splitDefinitions.map((def) => {
    const match = (regex.exec(def) as RegExpExecArray);
    return {
      definition: match[2],
      pronunciation: match[1],
    };
  });
}

function transformDefinitionFrom(
  definitions: backend.DictionaryDefinition[]) : string {
  return definitions.map(
    (def) => (def.pronunciation
      ? `[${def.pronunciation}] ${def.definition}`
      : def.definition),
  ).join('<br>');
}

function transformDefinitionFromPinyin(
  definitions: backend.DictionaryDefinition[]) : string[] {
  return definitions
    .map((def) => def.pronunciation)
    .filter((pro) => pro) as string[];
}

function transformTo(fields : backend.Fields) : FrontFields {
  const frontFields : FrontFields = {
    word: fields.word,
    sentence: {
      sentence: fields.sentence,
      source: undefined,
    },
    englishDefn: transformDefinition(fields.englishDefn),
    chineseDefn: transformDefinition(fields.chineseDefn),
    pinyin: '',
    images: fields.images,
  };
  return frontFields;
}

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
      stepsState: {} as StateMap,
      originalValues: ({} as FrontFields),
      newValues: ({} as FrontFields),
      note: new backend.RawAnkiNote(),
    };
  },
  getters: {

    word: (state) => {
      return state.newValues.word || state.originalValues.word;
    },
    sentence: (state) => {
      const sentence =
        state.newValues.sentence || state.originalValues.sentence;
      return sentence?.sentence;
    },
    sentenceSource: (state) => {
      const sentence =
        state.newValues.sentence || state.originalValues.sentence;
      return sentence?.source;
    },
    englishDefn: (state) => {
      return state.newValues.englishDefn ||
        state.originalValues.englishDefn || [];
    },
    chineseDefn: (state) => {
      return state.newValues.chineseDefn ||
        state.originalValues.chineseDefn || [];
    },
    pinyin: (state) => {
      return state.newValues.pinyin || state.originalValues.pinyin;
    },
    images: (state) => {
      return state.newValues.images || state.originalValues.images;
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
      this.originalValues = transformTo(this.note.fields);
      this.newValues = { word: this.originalValues.word };
      this.newValues.word = this.originalValues.word;

      this.currentStep = StepsEnum.SENTENCE;
      this.currentStepIndex = 0;
      this.flow = true;
    },

    getChanged() : backend.Fields {
      const fields = backend.Fields.createFrom({});
      fields.word = this.word;
      if (this.newValues.sentence) {
        fields.sentence = this.newValues.sentence.sentence;
      }
      if (this.newValues.englishDefn) {
        fields.englishDefn = transformDefinitionFrom(
          this.newValues.englishDefn);
      }
      if (this.newValues.chineseDefn) {
        fields.chineseDefn = transformDefinitionFrom(
          this.newValues.chineseDefn);
      }
      if (this.newValues.chineseDefn || this.newValues.englishDefn) {
        // If either of these has been set. Recalculate pinyin
        const oldPinyin = (this.originalValues.pinyin || '')
          .split(',')
          .map(pinyin => pinyin.trim())
          .filter(pinyin => pinyin.length > 0);

        fields.pinyin = [...new Set([
          ...transformDefinitionFromPinyin(this.chineseDefn),
          ...transformDefinitionFromPinyin(this.englishDefn),
          ...oldPinyin,
        ])].join(', ');
      }

      if (this.newValues.images) {
        fields.images = this.newValues.images;
      }

      return fields;
    },

    updateSentence(sentence: backend.Sentence) {
      this.newValues.sentence = sentence;
      this.stepsState[StepsEnum.SENTENCE] = StepState.PREVIEW;
    },

    updateDefinition(
      definitions: backend.DictionaryDefinition[],
      defType: string,
    ) {
      // const definitions = newDefinitions.map(
      //   (def) => `[${def.pronunciation}] ${def.definition}`,
      // ).join('<br>');
      if (defType === 'english') {
        this.newValues.englishDefn = definitions;
        this.stepsState[StepsEnum.ENGLISH] = StepState.PREVIEW;
      } else {
        this.newValues.chineseDefn = definitions;
        this.stepsState[StepsEnum.CHINESE] = StepState.PREVIEW;
      }
      // let pinyin = new Set();
      // if (this.newValues.pinyin !== undefined) {
      //   pinyin = new Set(this.newValues.pinyin.split(', '));
      // }
      // newDefinitions.forEach((def) => {
      //   const pronunciation = def.pronunciation.replace(/\s/g, '');
      //   pinyin.add(pronunciation);
      // });
      // pinyin.delete('');
      // this.newValues.pinyin = [...pinyin].join(', ');
    },

    updateImages(newImages: backend.ImageInfo[]) {
      this.newValues.images = newImages;
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
