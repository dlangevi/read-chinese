import { performance } from 'perf_hooks';
import fs from 'fs';
import { isKnown, knownArray } from './knownWords';
import { loadJieba } from './segmentation';
import { getDefinition } from './dictionaries';
import {
  dbGetBooks, dbGetBookById, dbAddBook, dbBookExists, dbSaveWordTable,
  dbGetBook, dbLoadWordTable, knex,
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

async function loadBook(bookId) {
  const book = await dbGetBookById(bookId);
  book.wordTable = await dbLoadWordTable(book);
  computeStats(book);
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

function computeStats(book) {
  let totalKnownWords = 0;
  let totalWords = 0;
  Object.entries(book.wordTable).forEach(([word, frequency]) => {
    totalWords += frequency;
    if (isKnown(word)) {
      totalKnownWords += frequency;
    }
  });
  book.totalKnownWords = totalKnownWords;
  book.totalWords = totalWords;
}

// This is where I get tripped up on the seperation layer. This is a db
// specific operation
export async function topWords() {
  const top = await knex('frequency')
    .select('word')
    .sum({ occurance: 'count' })
    .whereNotIn('word', knownArray())
    .groupBy('word')
    .orderBy('occurance', 'desc')
    .limit(200);

  return top.map((row) => {
    row.definition = getDefinition(row.word);
    return row;
  });
}

export function initLibraryIpc(ipcMain) {
  ipcMain.handle('loadBooks', async () => {
    const books = await dbGetBooks();
    books.forEach((book) => {
      const img = fs.readFileSync(book.cover).toString('base64');
      book.imgData = img;
    });
    return books;
  });
  ipcMain.handle('learningTarget', async () => {
    const start = performance.now();
    const words = await topWords();
    const end = performance.now();
    console.log(`Learning target took ${(end - start) / 1000}s`);
    return words;
  });

  ipcMain.handle('loadBook', async (event, bookId) => {
    const book = await loadBook(bookId);
    const img = fs.readFileSync(book.cover).toString('base64');
    book.imgData = img;

    return book;
  });
}
