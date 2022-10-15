import { IpcTypes } from '@/shared/ipcLoader';

export {};

declare global {
  interface Window {
    ipc: IpcTypes;
  }
}
