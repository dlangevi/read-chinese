import { contextBridge, ipcRenderer } from 'electron';

// Any strings in the argument list below will be transformed into a ipc
// function call.
// Is this a little hard to understand?
// Maybe.
// Was it fun to write and will cut down on the amount of boilerplate?
// Definitly
contextBridge.exposeInMainWorld(
  'ipc',
  ((...funs) => funs.reduce((acc, name) => {
    acc[name] = (...args) => ipcRenderer.invoke(name, ...args);
    return acc;
  }, {}))(
    'loadBooks',
    'loadBook',
    'learningTarget',
    'loadFlaggedCards',
    'getSentencesForWord',
    'getDefinitionsForWord',
    'getAnkiCard',
    'getAnkiNote',
    'updateAnkiCard',
    'addWord',
  ),
);
