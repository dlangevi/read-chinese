import fs from 'fs';
import {
  dbGetBooks, dbGetBookByID, dbAddBook, dbBookExists,
} from './database';

export function getBooks() {
  return getBooks();
}

export function addBook(author, title, cover, filepath) {
  dbAddBook(author, title, cover, filepath);
}
export function bookExists(author, title) {
  dbBookExists(author, title);
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

  ipcMain.handle('loadBook', (event, bookID) => {
    const book = dbGetBookByID(bookID);
    const img = fs.readFileSync(book.cover).toString('base64');
    book.imgData = img;

    return book;
  });
}
