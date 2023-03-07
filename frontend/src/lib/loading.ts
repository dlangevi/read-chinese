import { ref, inject } from 'vue';
import type { InjectionKey } from 'vue';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export class LoadingApi {
  displayText = ref('loading ...');
  active = ref(false);
  progressText = ref('');
  progressTotalSteps = ref(1);
  progressCurrentStep = ref(0);

  constructor() {
    EventsOn('setupProgress', (message: string, steps: number) => {
      this.progressText.value = message;
      this.progressTotalSteps.value = steps;
    });
    EventsOn('progress', () => {
      this.progressCurrentStep.value += 1;
    });
  }

  get progressMessage() {
    return this.progressText.value;
  }

  get progressPercent() {
    const steps = this.progressTotalSteps.value;
    const current = this.progressCurrentStep.value;
    return (current * 100 / steps).toFixed(0);
  }

  get progressSteps() {
    return this.progressTotalSteps.value;
  }

  get progressCurrent() {
    return this.progressCurrentStep.value;
  }

  get shouldShow() {
    return this.active.value;
  }

  get loadingText() {
    return this.progressText.value;
  }

  async withLoader(func:() => Promise<void | Error>) {
    // These will be set the the backend later
    this.progressText.value = '';
    this.progressTotalSteps.value = 1;
    this.progressCurrentStep.value = 0;

    this.active.value = true;
    const ret = await func();
    this.active.value = false;
    return ret;
  }
}

export const LoadingApiKey = Symbol('l') as InjectionKey<LoadingApi>;

export function useLoader():LoadingApi {
  return inject(LoadingApiKey) as LoadingApi;
}
