import { performance } from 'perf_hooks';
import fs from 'fs';
import { knownArray, isKnown, isKnownChar } from './knownWords';
import { loadJieba } from './segmentation';
import { getDefinition } from './dictionaries';
import {
  dbGetBooks, dbGetBookById, dbAddBook, dbBookExists, dbSaveWordTable,
  dbGetBook, knex,
} from './database';

export async function getBooks() {
  return dbGetBooks();
}

export async function addBook(author, title, cover, filepath) {
  await dbAddBook(author, title, cover, filepath);
  const book = await dbGetBook(author, title);
  const wordTable = await computeWordTable(book);
  return dbSaveWordTable(book, wordTable);
}
export async function bookExists(author, title) {
  return dbBookExists(author, title);
}

async function computeBookData(book) {
  // compute at runtime stuff I dont want to save right now
  await computeStats(book);
  const imgData = await fs.promises.readFile(book.cover);
  book.imgData = imgData.toString('base64');
}

async function computeExtraData(book) {
  const top = await knex('frequency')
    .select('word')
    .sum({ occurance: 'count' })
    .where('book', book.bookId)
    .groupBy('word');

  let probablyKnownWords = 0;
  let knownCharacters = 0;
  let totalCharacters = 0;
  top.forEach(({ word, occurance }) => {
    totalCharacters += word.length * occurance;
    let allKnown = true;
    Array.from(word).forEach((char) => {
      if (isKnownChar(char)) {
        knownCharacters += occurance;
      } else {
        allKnown = false;
      }
    });
    if (isKnown(word) || allKnown) {
      probablyKnownWords += occurance;
    }
  });

  book.probablyKnownWords = probablyKnownWords;
  book.knownCharacters = knownCharacters;
  book.totalCharacters = totalCharacters;
}

async function computeWordTargets(book) {
  const top = await knex('frequency')
    .select('word')
    .select('count')
    .where('book', book.bookId)
    .whereNotExists(function wordTable() {
      this.select('word')
        .from('words')
        .whereRaw('words.word==frequency.word');
    })
    .orderBy('count', 'desc');

  const targets = [
    80, 84, 86, 90, 92, 94, 96, 98, 100,
  ];
  const targetOccurances = targets.map(
    (target) => (target / 100) * book.totalWords,
  );
  const needToKnow = targetOccurances.map(
    (targetOccurance) => {
      let soFar = book.totalKnownWords;
      let needToLearn = 0;
      // I actually do need a loop here so I can short circut
      for (const entry of top) { // eslint-disable-line no-restricted-syntax
        if (soFar > targetOccurance) {
          break;
        }
        soFar += entry.count;
        needToLearn += 1;
      }
      return needToLearn;
    },
  );
  book.targets = targets;
  book.targetOccuances = targetOccurances;
  book.needToKnow = needToKnow;
}

async function loadBook(bookId) {
  const book = await dbGetBookById(bookId);
  await computeBookData(book);
  await computeExtraData(book);
  await computeWordTargets(book);
  return book;
}

async function computeWordTable(book) {
  console.log(`computing wordtable for ${book.filepath}`);
  const segText = await loadJieba(book.filepath);
  const wordTable = {};
  segText.forEach((sentence) => {
    sentence.forEach(([word, type]) => {
      if (type !== 3) return;
      if (word in wordTable) {
        wordTable[word] += 1;
      } else {
        wordTable[word] = 1;
      }
    });
  });
  return wordTable;
}

async function computeStats(book) {
  book.totalKnownWords = await knownWords(book);
  book.totalWords = await allWords(book);
}

// This is where I get tripped up on the seperation layer. This is a db
// specific operation
export async function topWords(bookIds) {
  const top = knex('frequency')
    .select('word')
    .sum({ occurance: 'count' })
    .whereNotIn('word', knownArray())
    .groupBy('word')
    .orderBy('occurance', 'desc')
    .limit(200);

  if (bookIds !== undefined && bookIds.length > 0) {
    top.whereIn('book', bookIds);
  }

  const results = await top;

  return results.map((row) => {
    row.definition = getDefinition(row.word);
    return row;
  });
}

export async function topUnknownWords(bookId, numWords) {
  const top = await knex('frequency')
    .select('word')
    .where('book', bookId)
    .whereNotExists(function wordTable() {
      this.select('word')
        .from('words')
        .whereRaw('words.word==frequency.word');
    })
    .orderBy('count', 'desc')
    .limit(numWords);

  return top;
}

async function knownWords(book) {
  const top = await knex('frequency')
    .sum({ occurance: 'count' })
    .where('book', book.bookId)
    .whereExists(function wordTable() {
      this.select('word')
        .from('words')
        .whereRaw('words.word==frequency.word');
    });
  return top[0].occurance;
}
async function allWords(book) {
  const top = await knex('frequency')
    .sum({ occurance: 'count' })
    .where('book', book.bookId);
  return top[0].occurance;
}

async function loadBooks() {
  const books = await dbGetBooks();
  console.time('loadBooks');
  await Promise.all(books.map((book) => computeBookData(book)));
  console.timeEnd('loadBooks');
  return books;
}
async function learningTarget(bookIds) {
  const start = performance.now();
  const words = await topWords(bookIds);
  const end = performance.now();
  console.log(`Learning target took ${(end - start) / 1000}s`);
  return words;
}

export const bookLibraryIpc = {
  loadBooks, learningTarget, loadBook, topUnknownWords,
};
