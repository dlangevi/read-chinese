import path from 'path';
import fs from 'fs';
import type {
  Book, UnknownWordEntry, HskLevel, HskVersion,
} from './types';
import {
  initBookStats,
} from './types';
import { isKnown, isKnownChar } from './knownWords';
import { loadJieba } from './segmentation';
import {
  dbGetBooks, dbGetBookById, dbAddBook, dbBookExists, dbSaveWordTable,
  dbGetBook, getKnex,
} from './database';

export async function getBooks(bookIds?: number[]) {
  return dbGetBooks(bookIds);
}

export async function addBook(
  author:string,
  title:string,
  cover:string,
  filepath:string,
) {
  const inserted = await dbAddBook(author, title, cover, filepath);
  if (inserted) {
    const book = await dbGetBook(author, title);
    const wordTable = await computeWordTable(book);
    await dbSaveWordTable(book, wordTable);
  }
}
export async function bookExists(author:string, title:string) {
  return dbBookExists(author, title);
}

async function computeBookData(book:Book) {
  // compute at runtime stuff I dont want to save right now
  await computeStats(book);
}

async function computeExtraData(book:Book) {
  return getKnex()<{ word:string, occurance:number }[]>('frequency')
    .select('word')
    .sum({ occurance: 'count' })
    .where('book', book.bookId)
    .groupBy('word')
    .then((rows) => {
      let probablyKnownWords = 0;
      let knownCharacters = 0;
      let totalCharacters = 0;
      rows.forEach(({ word, occurance }) => {
        totalCharacters += word.length * occurance;
        let allKnown = true;
        const charArray:string[] = Array.from(word);
        charArray.forEach((char:string) => {
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

      book.stats.probablyKnownWords = probablyKnownWords;
      book.stats.knownCharacters = knownCharacters;
      book.stats.totalCharacters = totalCharacters;
    });
}

async function computeWordTargets(book:Book) {
  const top = await getKnex()('frequency')
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
    (target) => (target / 100) * book.stats.totalWords,
  );
  const needToKnow = targetOccurances.map(
    (targetOccurance) => {
      let soFar = book.stats.totalKnownWords;
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
  book.stats.targets = targets;
  book.stats.targetOccurances = targetOccurances;
  book.stats.needToKnow = needToKnow;
}

async function loadBook(bookId:number) {
  const book = await dbGetBookById(bookId);
  book.stats = initBookStats();
  await computeBookData(book);
  await computeExtraData(book);
  await computeWordTargets(book);
  return book;
}

async function deleteBook(bookId:number) {
  await getKnex()('books').where('bookId', bookId).del();
  await getKnex()('frequency').where('book', bookId).del();
}

async function computeWordTable(book:Book) {
  console.log(`computing wordtable for ${book.filepath}`);
  const segText = await loadJieba(book);
  const wordTable:{
    [key:string]: number;
  } = {};
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

async function computeStats(book:Book) {
  book.stats.totalKnownWords = await knownWords(book);
  book.stats.totalWords = await allWords(book);
}

// This is where I get tripped up on the seperation layer. This is a db
// specific operation
export async function learningTarget(bookIds?:number[])
  : Promise<UnknownWordEntry[]> {
  const top = getKnex()<{ word:string, occurance:number }[]>('frequency')
    .select('word')
    .sum({ occurance: 'count' })
    .whereNotExists(function wordTable() {
      this.select('word')
        .from('words')
        .whereRaw('words.word==frequency.word');
    })
    .groupBy('word')
    .orderBy('occurance', 'desc')
    .limit(200);

  if (bookIds !== undefined && bookIds.length > 0) {
    top.whereIn('book', bookIds);
  }

  return top;
}

async function hskWords(version:HskVersion, level:HskLevel)
  : Promise<UnknownWordEntry[]> {
  const hskPath = path.join(
    __dirname,
    '../assets/HSK/',
    version,
    `L${level}.txt`,
  );
  const txt = fs.readFileSync(hskPath, {
    encoding: 'utf-8',
    flag: 'r',
  });
  const words = txt.split('\n');
  return words
    .map((word) => word.trim())
    .filter((word:string) => !isKnown(word))
    .filter((word:string) => word.length > 0)
    .map((word:string) => ({
      word,
    }));
}

export async function topUnknownWords(bookId:number, numWords:number) {
  const top = await getKnex()('frequency')
    .select('word')
    .where('book', bookId)
    .whereNotExists(function wordTable() {
      this.select('word')
        .from('words')
        .whereRaw('words.word==frequency.word');
    })
    .orderBy('count', 'desc')
    .limit(numWords);

  return top.map(({ word }) => word);
}

async function knownWords(book:Book) {
  const top = await getKnex()('frequency')
    .sum({ occurance: 'count' })
    .where('book', book.bookId)
    .whereExists(function wordTable() {
      this.select('word')
        .from('words')
        .whereRaw('words.word==frequency.word');
    });
  return top[0].occurance;
}
async function allWords(book:Book) {
  const top = await getKnex()('frequency')
    .sum({ occurance: 'count' })
    .where('book', book.bookId);
  return top[0].occurance;
}

async function loadBooks() {
  const books = await dbGetBooks();
  await Promise.all(books.map((book) => {
    book.stats = initBookStats();
    return computeBookData(book);
  }));
  return books;
}

async function setFavorite(bookId:number, isFavorite:boolean) {
  return getKnex()('books').where('bookId', bookId).update({
    favorite: isFavorite,
  });
}

async function setRead(bookId:number, hasRead:boolean) {
  return getKnex()('books').where('bookId', bookId).update({
    has_read: hasRead,
  });
}

async function totalRead() {
  const top = await getKnex()('frequency')
    .sum({ totalWords: 'count' })
    .whereExists(function wordTable() {
      this.select('bookId')
        .from('books')
        .where('has_read', true)
        .whereRaw('books.bookId==frequency.book');
    });
  return top[0];
}

export const bookLibraryIpc = {
  loadBooks,
  learningTarget,
  loadBook,
  topUnknownWords,
  deleteBook,
  setFavorite,
  setRead,
  totalRead,
  hskWords,
};
