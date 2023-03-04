import { defineStore } from 'pinia';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
import { toRaw } from 'vue';
import {
  GetAnkiNote,
} from '@wailsjs/backend/ankiInterface';
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

export type LoadOptions = {
  word: string
  // If set load the prev values from source card
  // Either can pass the noteId or the already loaded note
  // ( for testing )
  sourceCardId?: number,
  sourceCard?: backend.RawAnkiNote,
  // a list of steps that are 'important' (if empty its all)
  keySteps?: StepsEnum[],
  includeChinese?: boolean
  hasImageApi?: boolean
  disableFlow?: boolean

}

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
    word: fields.word || '',
    sentence: {
      sentence: fields.sentence || '',
      source: undefined,
    },
    englishDefn: transformDefinition(fields.englishDefn || ''),
    chineseDefn: transformDefinition(fields.chineseDefn || ''),
    pinyin: fields.pinyin,
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
      word: '',
      noteId: (undefined as number | undefined),
      stepsState: {} as StateMap,
      originalValues: (undefined as FrontFields | undefined),
      newValues: ({} as FrontFields),
    };
  },
  getters: {

    sentence: (state) => {
      const sentence =
        state.newValues.sentence || state.originalValues?.sentence;
      return sentence?.sentence;
    },
    sentenceSource: (state) => {
      const sentence =
        state.newValues.sentence || state.originalValues?.sentence;
      return sentence?.source;
    },
    englishDefn: (state) => {
      return state.newValues.englishDefn ||
        state.originalValues?.englishDefn || [];
    },
    chineseDefn: (state) => {
      return state.newValues.chineseDefn ||
        state.originalValues?.chineseDefn || [];
    },
    pinyin: (state) => {
      return state.newValues.pinyin || state.originalValues?.pinyin;
    },
    images: (state) => {
      return state.newValues.images || state.originalValues?.images;
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
    async loadCard(options : LoadOptions) {
      this.flow = false;
      this.currentStep = StepsEnum.NONE;
      this.steps = [
        StepsEnum.SENTENCE,
        StepsEnum.ENGLISH,
        // toggle this based on user settings
        ...(options.includeChinese ? [StepsEnum.CHINESE] : []),
        ...(options.hasImageApi ? [StepsEnum.IMAGE] : []),
      ];
      this.steps.forEach((step) => {
        this.stepsState[step] = StepState.EMPTY;
      });

      if (options.sourceCard) {
        this.originalValues = transformTo(options.sourceCard.fields);
      } else if (options.sourceCardId) {
        const noteData = await GetAnkiNote(options.sourceCardId);
        this.noteId = options.sourceCardId;
        this.originalValues = transformTo(noteData.fields);
      } else {
        this.noteId = undefined;
        this.originalValues = undefined;
      }
      this.word = options.word;
      this.newValues = { word: options.word };

      if (options.keySteps) {
        const keySteps = options.keySteps;
        const firstStep = this.steps.find(
          (step) => keySteps.includes(step));
        if (!firstStep) {
          throw new Error('Somehow keySteps did not have any valid steps');
        }
        this.currentStep = firstStep;
      } else {
        this.currentStep = StepsEnum.SENTENCE;
      }
      this.currentStepIndex = this.steps.indexOf(this.currentStep);
      this.flow = !options.disableFlow;
    },

    getChanged() : backend.Fields {
      const fields : {
        word?: string,
        sentence?: string,
        englishDefn?: string,
        chineseDefn?: string,
        pinyin?: string,
        images?: backend.ImageInfo[],
      } = {}; // backend.Fields.createFrom({});
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
      // If either of the definitions has changed
      if (this.newValues.chineseDefn ||
          this.newValues.englishDefn ||
            // Or if the original pinyin was not set at all
            !this.originalValues?.pinyin) {
        // Then generate a new pinyin field
        const oldPinyin = (this.originalValues?.pinyin || '')
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
        fields.images = toRaw(this.newValues.images);
      }

      return backend.Fields.createFrom(fields);
    },

    updateSentence(sentence: backend.Sentence) {
      this.newValues.sentence = sentence;
      this.stepsState[StepsEnum.SENTENCE] = StepState.PREVIEW;
    },

    updateDefinition(
      definitions: backend.DictionaryDefinition[],
      defType: string,
    ) {
      if (defType === 'english') {
        this.newValues.englishDefn = definitions;
        this.stepsState[StepsEnum.ENGLISH] = StepState.PREVIEW;
      } else {
        this.newValues.chineseDefn = definitions;
        this.stepsState[StepsEnum.CHINESE] = StepState.PREVIEW;
      }
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
      if (this.currentStepIndex + 1 === this.steps.length) {
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

export function formatDefinition(definition : backend.DictionaryDefinition) {
  if (definition.pronunciation) {
    return `[${definition.pronunciation}] ${definition.definition}`;
  }
  return definition.definition;
}

export function getImageSrc(image : backend.ImageInfo | undefined) : string {
  if (image === undefined) {
    return '';
  } else if (image.url !== undefined) {
    return image.url;
  } else if (image.imageData !== undefined) {
    return `data:image/png;base64, ${image.imageData}`;
  } else {
    return '';
  }
}
