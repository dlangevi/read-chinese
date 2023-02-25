import { ref, inject } from 'vue';
import type { InjectionKey } from 'vue';

export class LoadingApi {
  displayText = ref('loading ...');
  active = ref(false);

  get shouldShow() {
    return this.active.value;
  }

  get loadingText() {
    return this.displayText.value;
  }

  async withLoader(func:() => Promise<any>, message:string) {
    this.displayText.value = message;
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
