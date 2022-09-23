import {
  ipcMain, contextBridge, ipcRenderer, dialog,
} from 'electron';
import { knownWordsIpc } from './background/knownWords';
import { bookLibraryIpc } from './background/bookLibrary';
import { ankiInterfaceIpc } from './background/ankiInterface';
import { dictionariesIpc } from './background/dictionaries';
import { generateSentencesIpc } from './background/generateSentences';
import { imageSearchIpc } from './background/imageSearch';
import { databaseIpc } from './background/database';

// Put this here for now
function filePicker(extension) {
  const file = dialog.showOpenDialogSync({
    properties: ['openFile'],
    filters: [
      { name: 'Any File', extensions: [extension] },
    ],
  });
  return file[0];
}

// Concatenate all the functions we want accessed via IPC
// TODO should have some way to ensure no name collisions
const ipcFunctions = {
  ...bookLibraryIpc,
  ...knownWordsIpc,
  ...generateSentencesIpc,
  ...ankiInterfaceIpc,
  ...dictionariesIpc,
  ...databaseIpc,
  ...imageSearchIpc,
  filePicker,
};

// To be called from background.js to initialize handlers
export function initIpcMain() {
  Object.entries(ipcFunctions).forEach(([name, fn]) => {
    // For each function register it to be handled and drop the
    // event argument (here as _)
    ipcMain.handle(name, (_, ...args) => fn(...args));
  });
}

// To be called from preload to initialize callers
export function initIpcRenderer() {
  contextBridge.exposeInMainWorld(
    'ipc',
    // Map the functions to a object with keys that allow the calling of each
    // function by its name in vue land
    Object.keys(ipcFunctions).reduce((acc, name) => {
      acc[name] = (...a) => ipcRenderer.invoke(name, ...a);
      return acc;
    }, {}),
  );
}
