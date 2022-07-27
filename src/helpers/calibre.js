// import Database from 'better-sqlite3';

// Load a users calibre database
// const db = new Database('/home/dlangevi/chinese/chinese books/metadatabk.db',
// {verbose: console.log});
//
import { ipcMain } from 'electron';
import { Calibre } from 'node-calibre';
import fs from 'fs';
import { bookExists, addBook, getBooks } from './database';

// Create Calibre instance
const calibre = new Calibre({ library: '/home/dlangevi/chinese/chinese books/' });

async function getCalibreBooks() {
  const result = await calibre.run('calibredb list', {
    // limit: 10,
    forMachine: null,
    fields: 'cover,authors,title,formats',
  });
  const books = JSON.parse(result);
  return books;
}

export async function importCalibreBooks() {
  const books = await getCalibreBooks();
  Object.values(books).forEach((book) => {
    // Overwrite everythin every time for now
    if (bookExists(book.authors, book.title)) {
      // For now just add the ones that already have txt files.
      const txtBooks = book.formats.filter((path) => path.match(/.*.txt/));
      if (txtBooks.length > 0) {
        addBook(book.authors, book.title, book.cover, txtBooks[0]);
      }
    }
  });
}

export function initLibraryIpc() {
  console.log('setting up lib ipc');
  ipcMain.on('need-books', async (event) => {
    console.log('got request for books');
    const books = getBooks();
    console.log(books);
    books.forEach((book) => {
      const img = fs.readFileSync(book.cover).toString('base64');
      // Haha this gets around no reasign warning...
      const fake = book;
      fake.imgData = img;
    });
    event.reply('give-books', books);
  });
}
