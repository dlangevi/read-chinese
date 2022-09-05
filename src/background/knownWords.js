import { updateWord, wordExists } from './database';
// In the future if we need performance we could have a mirror of the sql data
// as an inmemory object but for now will just have these be wrapper calls
// around db operations

// For now the db code will update the word set here on each addition.
// In the future there should not be two seperate sets of words
export function addWord(word, age = 0, hasFlashCard = false) {
  updateWord(word, age, hasFlashCard);
}

export function saveLegacyWords(words) {
  Object.entries(words).forEach(([word, entry]) => {
    console.log(`Inserting ${word}`);
    addWord(word, entry.interval);
  });
}

export function isKnown(word) {
  return wordExists(word);
}

export function initWordsIpc(ipcMain) {
  ipcMain.handle('addWord', (event, word) => {
    // These will be from markedLearned so stick to prior convention for now
    addWord(word, 10000);
  });
}
