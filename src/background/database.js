// Dont know if this is a silly idea, but want to isolate all the actual reads
// and writes in a seperate file. So anything that is saved or read from
// persistant storage will have to go through functions prefixed 'db' in this
// file. If we want to swap out backends later on at least all the code to be
// changed will be in the same place
import Store from 'electron-store';
import Knex from 'knex';
// For now we do the sync whenever the db changes.
import knexConfigMap from '../../knexfile';

const knexConfig = knexConfigMap[process.env.NODE_ENV];
const knex = Knex(knexConfig);

// This is called and awaited before before anyother code can run
export async function initializeDatabase() {
  await knex.migrate.latest(knexConfig).catch((err) => {
    console.log(err);
  });
}

// Books and metadata can be stored in electron-store for now since they should
// be low footprint
const bookStore = new Store({ name: 'books' });
const metadataStore = new Store({ name: 'metadata' });

/** *********************************
 *
 * Metadata
 *
 ********************************** */

export function updateTimesRan() {
  const timesRan = metadataStore.get('ran', 0);
  metadataStore.set('ran', timesRan + 1);
}

export function getTimesRan() {
  return metadataStore.get('ran', 0);
}

/** *********************************
 *
 * Known Words + Flash Cards
 *
 * eachRow: {
 *    word: string,
 *    has_flash_card: boolean,
 *    has_sentence: boolean,
 *    interval: integer, // Anki flashcard interval
 * }
 *
 ********************************** */

// Adds if not exists
export async function dbUpdateWord(word, interval = 0, hasFlashCard = false) {
  const exists = await knex('words').select().where('word', word);
  if (exists.length === 0) {
    console.log(`Adding new word: ${word}`);
    knex('words')
      .insert({
        word,
        interval,
        has_flash_card: hasFlashCard,
      }).catch((err) => { console.log(err); });
  } else {
    knex('words')
      .where('word', word)
      .update({
        has_flash_card: hasFlashCard,
        interval,
      }).catch((err) => { console.log(err); });
  }
}

export async function dbLoadWords() {
  const rows = await knex('words')
    .select({ id: 'id', word: 'word' })
    .catch((error) => { console.log(error); });
  const words = new Set();
  rows.forEach((row) => {
    words.add(row.word);
  });
  return words;
}

/** *********************************
 *
 * Books
 *
 * currently indexed by combination of author and title
 *
 * bookKey: {
 *  author: string,
 *  title: string,
 *  txtFile: string, // path of where book txt file is stored
 *  cover: string, // path of where book cover image is stored
 *  bookID: string,
 *  wordTable: { word => count } count of each word in the book
 * }
 *
 ********************************** */
function bookKey(author, title) {
  return `${author}-${title}`;
}
export function dbAddBook(author, title, cover, filepath) {
  // For now just point to the actual txt file location in calibre. Later we will make our own copy
  const books = bookStore.get('booklist', {});
  books[bookKey(author, title)] = {
    author,
    title,
    txtFile: filepath,
    cover,
    bookID: bookKey(author, title),
  };
  bookStore.set('booklist', books);
}

export function dbGetBooks() {
  return Object.values(bookStore.get('booklist', {}));
}

export function dbGetBook(author, title) {
  return bookStore.get('booklist')[bookKey(author, title)];
}

export function dbGetBookByID(bookID) {
  return bookStore.get('booklist')[bookID];
}

// For now we will use author and title to do book uniqueness
export function dbBookExists(author, title) {
  const books = bookStore.get('booklist', {});
  return (bookKey(author, title) in books);
}
