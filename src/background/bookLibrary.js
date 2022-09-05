import fs from 'fs';
import { isKnown } from './knownWords';
import { loadJieba } from './segmentation';
import {
  dbGetBooks, dbGetBookByID, dbAddBook, dbBookExists,
} from './database';

export function getBooks() {
  return dbGetBooks();
}

export function addBook(author, title, cover, filepath) {
  dbAddBook(author, title, cover, filepath);
}
export function bookExists(author, title) {
  dbBookExists(author, title);
}

async function loadBook(bookID) {
  const book = dbGetBookByID(bookID);
  if (!book.wordTable) {
    book.wordTable = await computeWordTable(book);
    // Save wordTable
  }
  computeStats(book);
  return book;
}

async function computeWordTable(book) {
  console.log(`computing wordtable for ${book.txtFile}`);
  const segText = await loadJieba(book.txtFile);
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

export function initLibraryIpc(ipcMain) {
  ipcMain.handle('loadBooks', () => {
    const books = dbGetBooks();
    books.forEach((book) => {
      const img = fs.readFileSync(book.cover).toString('base64');
      book.imgData = img;
    });
    return books;
  });

  ipcMain.handle('loadBook', async (event, bookID) => {
    const book = await loadBook(bookID);
    const img = fs.readFileSync(book.cover).toString('base64');
    book.imgData = img;

    return book;
  });
}
