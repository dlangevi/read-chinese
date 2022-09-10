import { ipcMain, contextBridge, ipcRenderer } from 'electron';
import { knownWordsIpc } from './background/knownWords';
import { bookLibraryIpc } from './background/bookLibrary';
import { ankiInterfaceIpc } from './background/ankiInterface';
import { dictionariesIpc } from './background/dictionaries';
import { generateSentencesIpc } from './background/generateSentences';

// Concatenate all the functions we want accessed via IPC
// TODO should have some way to ensure no name collisions
const ipcFunctions = [
  ...bookLibraryIpc,
  ...knownWordsIpc,
  ...generateSentencesIpc,
  ...ankiInterfaceIpc,
  ...dictionariesIpc,
];

// To be called from background.js to initialize handlers
export function initIpcMain() {
  ipcFunctions.forEach((fn) => {
    // For each function register it to be handled and drop the
    // event argument (here as _)
    ipcMain.handle(fn.name, (_, ...args) => fn(...args));
  });
}

// To be called from preload to initialize callers
export function initIpcRenderer() {
  contextBridge.exposeInMainWorld(
    'ipc',
    // Map the functions to a object with keys that allow the calling of each
    // function by its name in vue land
    ipcFunctions.reduce((acc, fn) => {
      acc[fn.name] = (...a) => ipcRenderer.invoke(fn.name, ...a);
      return acc;
    }, {}),
  );
}
