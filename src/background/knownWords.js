import { dialog } from 'electron';
import fs from 'fs';
import { dbLoadWords, dbUpdateWord, dbUpdateWords } from './database';

// Memory cache of the set of known words for performance
let known = {};
let knownCharacters = new Set();
export async function syncWords() {
  known = await dbLoadWords();
  knownCharacters = new Set();
  known.forEach((word) => {
    Array.from(word).forEach((ch) => knownCharacters.add(ch));
  });
  console.log(`Known words: ${known.size}
Known characters: ${knownCharacters.size} `);
}

export function knownArray() {
  return [...known];
}

// For now the db code will update the word set here on each addition.
// In the future there should not be two seperate sets of words
export function addWord(word, age = 0, hasFlashCard = false) {
  known.add(word);
  dbUpdateWord(word, age, hasFlashCard);
}

// wordRows expects [{word, interval, has_flash_card}]
export function addWords(wordRows) {
  wordRows.forEach((row) => {
    known.add(row.word);
  });
  dbUpdateWords(wordRows);
}

export function importLegacyWords() {
  const wordsFile = dialog.showOpenDialogSync({
    properties: ['openFile'],
    filters: [
      { name: 'Json format', extensions: ['json'] },
    ],
  });
  if (wordsFile === undefined) {
    return;
  }
  const contents = fs.readFileSync(wordsFile[0], {
    encoding: 'utf-8',
    flags: 'r',
  });
  try {
    // Dont worry about validating since this is only for my personal use
    const words = JSON.parse(contents);
    const wordRows = Object.entries(words).map(
      ([word, entry]) => ({ word, interval: entry.interval }),
    );
    addWords(wordRows);
  } catch (error) {
    console.error(error);
  }
}

export function importCSVWords() {
  // TODO handle bad selections
  const wordsFile = dialog.showOpenDialogSync({
    properties: ['openFile'],
    filters: [
      { name: 'one per line', extensions: ['csv'] },
    ],
  });
  if (wordsFile === undefined) {
    return;
  }
  const contents = fs.readFileSync(wordsFile[0], {
    encoding: 'utf-8',
    flags: 'r',
  });
  // For now its just one word per line no other data?
  const words = contents.split('\n');
  words.forEach((word) => {
    addWord(word);
  });
}

export function isKnown(word) {
  return known.has(word);
}

export function isKnownChar(char) {
  return knownCharacters.has(char);
}

export const knownWordsIpc = [addWord, importLegacyWords, importCSVWords];
