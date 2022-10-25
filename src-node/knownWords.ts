import fs from 'fs';
import type { KnownWords } from './types';
import { isInDictionary } from './dictionaries';
import {
  dbLoadWords, dbUpdateWord, dbUpdateWords, getOptionValue,
} from './database';

// Memory cache of the set of known words for performance
let known: KnownWords = {};
let knownCharacters = new Set();
export async function syncWords() {
  await dbLoadWords().then((words) => {
    known = words;
  });

  knownCharacters = new Set();
  Object.keys(known).forEach((word) => {
    Array.from(word).forEach((ch) => knownCharacters.add(ch));
  });
  console.log(`Known words: ${Object.keys(known).length}
Known characters: ${knownCharacters.size} `);
}

function wordStats() {
  return {
    words: Object.keys(known).length,
    characters: knownCharacters.size,
    wack: checkWords(),
  };
}

export function checkWords() {
  let wackWords = 0;
  Object.keys(known).forEach((word) => {
    if (!isInDictionary(word)) {
      wackWords += 1;
    }
  });
  return wackWords;
}

// For now the db code will update the word set here on each addition.
// In the future there should not be two seperate sets of words
export function addWord(word:any, age = 0, hasFlashCard = false) {
  known[word] = { interval: age };
  dbUpdateWord(word, age, hasFlashCard);
}

// wordRows expects [{word, interval, has_flash_card}]
export function addWords(wordRows:any) {
  wordRows.forEach((row:any) => {
    known[row.word] = { interval: row.interval };
  });
  dbUpdateWords(wordRows);
}

export function importCSVWords() {
  // PORTOVER
  // TODO handle bad selections
  // const wordsFile = dialog.showOpenDialogSync({
  //   properties: ['openFile'],
  //   filters: [
  //     { name: 'one per line', extensions: ['csv'] },
  //   ],
  // });
  const wordsFile = undefined;
  if (wordsFile === undefined) {
    return;
  }
  const contents = fs.readFileSync(wordsFile[0], {
    encoding: 'utf-8',
    flag: 'r',
  });
  // For now its just one word per line no other data?
  const words = contents.split('\n');
  words.forEach((word) => {
    addWord(word);
  });
}

let knownInterval = 100;
export function updateInterval() {
  knownInterval = getOptionValue('KnownInterval', 100);
}

export function isWellKnown(word:string) {
  return (word in known) && known[word].interval >= knownInterval;
}

export function isKnown(word:string) {
  return (word in known);
}

export function isKnownChar(char:string) {
  return knownCharacters.has(char);
}

export const knownWordsIpc = {
  addWord,
  importCSVWords,
  wordStats,
};
