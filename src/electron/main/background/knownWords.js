import { dialog } from 'electron';
import fs from 'fs';
import { isInDictionary } from './dictionaries';
import {
  dbLoadWords, dbUpdateWord, dbUpdateWords, getOptionValue,
} from './database';

// Memory cache of the set of known words for performance
let known = {};
let knownCharacters = new Set();
export async function syncWords() {
  known = await dbLoadWords();
  knownCharacters = new Set();
  Object.keys(known).forEach((word) => {
    Array.from(word).forEach((ch) => knownCharacters.add(ch));
  });
  console.log(`Known words: ${Object.keys(known).length}
Known characters: ${knownCharacters.size} `);
}

export function checkWords() {
  let wackWords = 0;
  Object.keys(known).forEach((word) => {
    if (!isInDictionary(word)) {
      wackWords += 1;
      console.log(word);
    }
  });

  console.log(`You got ${wackWords} wack ass words`);
}

// For now the db code will update the word set here on each addition.
// In the future there should not be two seperate sets of words
export function addWord(word, age = 0, hasFlashCard = false) {
  known[word] = { interval: age };
  dbUpdateWord(word, age, hasFlashCard);
}

// wordRows expects [{word, interval, has_flash_card}]
export function addWords(wordRows) {
  wordRows.forEach((row) => {
    known[row.word] = { interval: row.interval };
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

let knownInterval = getOptionValue('KnownInterval', 100);
export function updateInterval() {
  knownInterval = getOptionValue('KnownInterval', 100);
}

export function isKnown(word) {
  return (word in known) && known[word].interval >= knownInterval;
}

export function isKnownChar(char) {
  return knownCharacters.has(char);
}

export const knownWordsIpc = { addWord, importLegacyWords, importCSVWords };
