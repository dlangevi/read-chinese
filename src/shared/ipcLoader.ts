import {
  ipcMain, contextBridge, ipcRenderer, dialog,
} from 'electron';
import { knownWordsIpc } from '@/main/knownWords';
import { bookLibraryIpc } from '@/main/bookLibrary';
import { ankiInterfaceIpc } from '@/main/ankiInterface';
import { dictionariesIpc } from '@/main/dictionaries';
import { generateSentencesIpc } from '@/main/generateSentences';
import { imageSearchIpc } from '@/main/imageSearch';
import { databaseIpc } from '@/main/database';
import { calibreIpc } from '@/main/calibre';

// Put this here for now
function filePicker(extension:string) {
  const file = dialog.showOpenDialogSync({
    properties: ['openFile'],
    filters: [
      { name: 'Any File', extensions: [extension] },
    ],
  });
  if (file === undefined) {
    return 'error';
  }
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
  ...calibreIpc,
  filePicker,
};

// When called as an ipc, the function on the renderer side will
// return a promise. So we massage the types of actual implementations
// to reflect that
type Promised<Type> =
  // If the type already has a Promise being returned
  Type extends (...params: infer P) => Promise<infer R>
    // We can just return that same promise
    ? (...params: P) => Promise<R>
    // Otherwise take the return value
    : Type extends (...params: infer P) => infer R
      // And wrap it in a promise
      ? (...params: P) => Promise<R>
      : never;

type BasicIpcTypes = typeof ipcFunctions;
export type IpcTypes = {
  [Property in keyof BasicIpcTypes]: Promised<BasicIpcTypes[Property]>
};

type IpcFunction = (...args: any[]) => Promise<any>;

// To be called from background.js to initialize handle
export function initIpcMain() {
  // Wipe out the types of ipcFunctions here because we know it will be typesafe
  Object.entries(ipcFunctions).forEach(([name, fn]: [string, any]) => {
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
    Object.keys(ipcFunctions).reduce((
      acc:{
        [fnName:string]: IpcFunction
      },
      name:string,
    ) => {
      acc[name] = (...a:any[]) => ipcRenderer.invoke(name, ...a);
      return acc;
    }, {}),
  );
}
