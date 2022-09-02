// import { ipcRenderer } from 'electron';

// window.ipcRenderer = ipcRenderer;
// window.dave = 32;

import { contextBridge, ipcRenderer } from 'electron';

const call = ipcRenderer.invoke;
contextBridge.exposeInMainWorld('ipc', {
  loadBooks: () => call('loadBooks'),
  loadBook: (title) => call('loadBook', title),
  learningTarget: () => call('learningTarget'),
  loadFlaggedCards: () => call('flaggedCards'),
  getSentencesForWord: (word) => call('getSentencesForWord', word),
  getAnkiCard: (word) => call('getAnkiCard', word),
  getAnkiNote: (word) => call('getAnkiNote', word),
  updateAnkiCard: (noteID, fields) => call('updateAnkiCard', noteID, fields),
  addWord: (word) => call('addWord', word),
});
