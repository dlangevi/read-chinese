// import wordStats from './wordStats.js';
// import config from './config.js';
import { loadWords, updateWord } from './database';

// @todo save and load this from the database,
// and handle per user word lists

let known = {};
let knownCharacters = new Set();
export async function syncWords() {
  known = await loadWords();
  knownCharacters = new Set();
  Object.keys(known).forEach((word) => {
    Array.from(word).forEach((ch) => knownCharacters.add(ch));
  });
  console.log(`Known words: ${Object.keys(known).length}
Known characters: ${knownCharacters.size} `);
}

/**
 * Format a date to the form 'YYYY-MM-DD'
 * @param {Date} date
 * @return {string}
 */
export function toDateString(date) {
  const year = date.getFullYear();
  let month = date.getMonth() + 1;
  if (month < 10) {
    month = `0${month}`;
  }
  let day = date.getDate();
  if (day < 10) {
    day = `0${day}`;
  }
  return `${year}-${month}-${day}`;
}

/**
 * Return the current date formated 'YYYY-MM-DD'
 * @return {string}
 */
function currentDateString() {
  return toDateString(new Date());
}

// For now the db code will update the word set here on each addition.
// In the future there should not be two seperate sets of words
export function addWord(word, age = 0, hasFlashCard = false) {
  // If this is a new word, add it with the current date
  if (!(word in known)) {
    known[word] = {
      added: currentDateString(),
      interval: age,
    };
    console.log(`Adding new word ${word} ${JSON.stringify(known[word])}`);
  } else if (known[word].interval !== age) {
    // else just update the interval
    known[word].interval = age;
    console.log(`Updating Word interval for ${word} to ${age}`);
  }
  updateWord(word, age, hasFlashCard);
}

export function saveLegacyWords(words) {
  Object.entries(words).forEach(([word, entry]) => {
    console.log(`Inserting ${word}`);
    addWord(word, entry.interval);
  });
}

export function isKnown(word, howKnown = 0) {
  // if word is completly unknown return false
  if (!(word in known)) {
    return false;
  }
  // we know it at least somewhat known
  return known[word].interval >= howKnown;
}

export function initWordsIpc(ipcMain) {
  ipcMain.handle('addWord', (event, word) => {
    // These will be from markedLearned so stick to prior convention for now
    addWord(word, 10000);
  });
}
