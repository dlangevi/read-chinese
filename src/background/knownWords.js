import { loadWords, updateWord } from './database';

// Memory cache of the set of known words for performance
let known = {};
let knownCharacters = new Set();
export async function syncWords() {
  known = await loadWords();
  knownCharacters = new Set();
  known.forEach((word) => {
    Array.from(word).forEach((ch) => knownCharacters.add(ch));
  });
  console.log(`Known words: ${known.size}
Known characters: ${knownCharacters.size} `);
}

// For now the db code will update the word set here on each addition.
// In the future there should not be two seperate sets of words
export function addWord(word, age = 0, hasFlashCard = false) {
  known.add(word);
  updateWord(word, age, hasFlashCard);
}

export function saveLegacyWords(words) {
  Object.entries(words).forEach(([word, entry]) => {
    console.log(`Inserting ${word}`);
    addWord(word, entry.interval);
  });
}

export function isKnown(word) {
  return known.has(word);
}

export function initWordsIpc(ipcMain) {
  ipcMain.handle('addWord', (event, word) => {
    // These will be from markedLearned so stick to prior convention for now
    addWord(word, 10000);
  });
}
