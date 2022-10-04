import {
  ipcMain, contextBridge, ipcRenderer, dialog,
} from 'electron';
import { knownWordsIpc } from '../main/knownWords';
import { bookLibraryIpc } from '../main/bookLibrary';
import { ankiInterfaceIpc } from '../main/ankiInterface';
import { dictionariesIpc } from '../main/dictionaries';
import { generateSentencesIpc } from '../main/generateSentences';
import { imageSearchIpc } from '../main/imageSearch';
import { databaseIpc } from '../main/database';

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
