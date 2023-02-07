import { ref, inject } from 'vue';
import type { InjectionKey } from 'vue';
// import type { Ref } from 'vue';
//

type ToastType = 'info' | 'success' | 'warning' | 'error'
export class MessageApi {
  showToast = ref(false);
  toastType = ref<ToastType>('info');
  text = ref('');

  // constructor() {
  // }

  get shouldShow() {
    return this.showToast.value;
  }

  get currentType() {
    return this.toastType.value;
  }

  get message() {
    return this.text.value;
  }

  info(text: string) {
    this.#showMessage(text, 'info');
  }

  success(text: string) {
    this.#showMessage(text, 'success');
  }

  error(text: string) {
    this.#showMessage(text, 'error', true);
  }

  warn(text: string) {
    this.#showMessage(text, 'warning');
  }

  #showMessage(text:string, type:ToastType, keepOpen?:boolean) {
    this.showToast.value = true;
    this.toastType.value = type;
    this.text.value = text;
    if (!keepOpen) {
      setTimeout(() => {
        this.hide();
      }, 3000);
    }
  }

  hide() {
    this.showToast.value = false;
  }
}

export const MessageApiKey = Symbol('m') as InjectionKey<MessageApi>;

export function useMessage():MessageApi {
  return inject(MessageApiKey) as MessageApi;
}
