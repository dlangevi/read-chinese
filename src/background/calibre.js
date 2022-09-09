// import Database from 'better-sqlite3';

// Load a users calibre database
// const db = new Database('/home/dlangevi/chinese/chinese books/metadatabk.db',
// {verbose: console.log});
//
import { Calibre } from 'node-calibre';
import { addBook } from './bookLibrary';

async function getCalibreBooks(calibreDir) {
  // Create Calibre instance
  const calibre = new Calibre({ library: calibreDir });
  const result = await calibre.run('calibredb list', {
    // limit: 10,
    forMachine: null,
    fields: 'cover,authors,title,formats',
  });
  const books = JSON.parse(result);
  return books;
}

export async function importCalibreBooks(calibreDir) {
  const books = await getCalibreBooks(calibreDir);
  Object.values(books).forEach((book) => {
    // Overwrite everythin every time for now
    // if (bookExists(book.authors, book.title)) {
    // For now just add the ones that already have txt files.
    // TODO do some conversion of epubs etc
    // (is epub-convert avaliable on all platforms?)
    const txtBooks = book.formats.filter((path) => path.match(/.*.txt/));
    if (txtBooks.length > 0) {
      addBook(book.authors, book.title, book.cover, txtBooks[0]);
    }
    // }
  });
}
