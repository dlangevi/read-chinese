// import { ipcRenderer } from 'electron';

// window.ipcRenderer = ipcRenderer;
// window.dave = 32;

import { contextBridge, ipcRenderer } from 'electron';

// Maybe later we will want to set this up but meh
// const validChannels = ['READ_FILE', 'WRITE_FILE'];
contextBridge.exposeInMainWorld('ipc', {
  loadBooks: () => ipcRenderer.invoke('loadBooks'),
  loadBook: (title) => ipcRenderer.invoke('loadBook', title),
  learningTarget: () => ipcRenderer.invoke('learningTarget'),
  loadFlaggedCards: () => ipcRenderer.invoke('flaggedCards'),
  getSentencesForWord: (word) => ipcRenderer.invoke('getSentencesForWord', word),
  getAnkiCard: (word) => ipcRenderer.invoke('getAnkiCard', word),
  getAnkiNote: (word) => ipcRenderer.invoke('getAnkiNote', word),
  updateAnkiCard: (noteID, fields) => ipcRenderer.invoke('updateAnkiCard', noteID, fields),
});
