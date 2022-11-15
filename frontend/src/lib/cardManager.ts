import type {
  backend,
} from '@wailsjs/models';
import { StepsEnum } from '@/components/CardCreatorSteps/StepsEnum';
// Should take an anki card, and a list of steps to be computed

export class CardManager {
  note : backend.RawAnkiNote;
  steps : StepsEnum[];
  constructor(ankiCard : backend.RawAnkiNote, steps: StepsEnum) {
    this.note = ankiCard;
    this.steps = [steps];
  }
}
