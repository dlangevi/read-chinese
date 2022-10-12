import { ipcTypes } from '../shared/ipcLoader';

export {};

declare global {
  interface Window {
    ipc: ipcTypes;
  }
}
