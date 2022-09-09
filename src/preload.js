import { contextBridge, ipcRenderer } from 'electron';

const call = ipcRenderer.invoke;
contextBridge.exposeInMainWorld('ipc', {
  loadBooks: () => { return call('loadBooks'); },
  loadBook: (title) => { return call('loadBook', title); },
  learningTarget: (books) => { return call('learningTarget', books); },
  loadFlaggedCards: () => { return call('flaggedCards'); },
  getSentencesForWord: (word) => { return call('getSentencesForWord', word); },
  getDefinitionsForWord: (word) => {
    return call('getDefinitionsForWord', word);
  },
  getAnkiCard: (word) => { return call('getAnkiCard', word); },
  getAnkiNote: (word) => { return call('getAnkiNote', word); },
  updateAnkiCard: (noteID, fields) => {
    return call('updateAnkiCard', noteID, fields);
  },
  addWord: (word) => { return call('addWord', word); },
});
