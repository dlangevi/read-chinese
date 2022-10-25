import type { Express } from 'express';
import { knownWordsIpc } from './knownWords';
import { bookLibraryIpc } from './bookLibrary';
import { ankiInterfaceIpc } from './ankiInterface';
import { dictionariesIpc } from './dictionaries';
import { generateSentencesIpc } from './generateSentences';
import { imageSearchIpc } from './imageSearch';
import { databaseIpc } from './database';
import { calibreIpc } from './calibre';

// Put this here for now
function filePicker(extension:string) {
  // const file = dialog.showOpenDialogSync({
  //   properties: ['openFile'],
  //   filters: [
  //     { name: 'Any File', extensions: [extension] },
  //   ],
  // });
  // if (file === undefined) {
  //   return 'error';
  // }
  // return file[0];
  console.log(`implement ${extension}`);
  return 'error';
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

export function registerIpc(app:Express) {
  app.post('/ipc', async (req, res) => {
    const functionName: keyof BasicIpcTypes = req.body.function;
    const args = JSON.parse(req.body.args);
    const ipcFunction:any = ipcFunctions[functionName];
    if (ipcFunction === undefined) {
      console.log({
        error: `${req.body.function}`,
      });
      res.json({
        error: `${req.body.function}`,
      });
    }
    const result = await ipcFunction(...args);
    res.json(result);
  });
}
