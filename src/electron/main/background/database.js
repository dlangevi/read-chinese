// Dont know if this is a silly idea, but want to isolate all the actual reads
// and writes in a seperate file. So anything that is saved or read from
// persistant storage will have to go through functions prefixed 'db' in this
// file. If we want to swap out backends later on at least all the code to be
// changed will be in the same place
import Store from 'electron-store';
import Knex from 'knex';
// For now we do the sync whenever the db changes.
import knexConfigMap from '../../../../knexfile.mjs';

console.log(knexConfigMap);
const knexConfig = knexConfigMap[import.meta.env.MODE];
export const knex = Knex(knexConfig);

// This is called and awaited before before anyother code can run
export async function initializeDatabase() {
  await knex.migrate.latest(knexConfig).catch((err) => {
    console.log(err);
  });
}

// Books and metadata can be stored in electron-store for now since they should
// be low footprint
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

// TODO, should there be some white list of valid keys?
export function getOptionValue(key, defaultValue) {
  return metadataStore.get(key, defaultValue);
}

function setOptionValue(key, value) {
  metadataStore.set(key, value);
}

export function dbSaveDict(name, path, type) {
  const dicts = metadataStore.get('dicts', {});
  dicts[name] = {
    path,
    type,
  };
  metadataStore.set('dicts', dicts);
}

export function dbDeleteDict(name) {
  const dicts = metadataStore.get('dicts', {});
  delete dicts[name];
  metadataStore.set('dicts', dicts);
}

export function dbLoadDicts() {
  return metadataStore.get('dicts', {});
}

export function dbSetPrimaryDict(dictName) {
  return metadataStore.set('primaryDict', dictName);
}
export function dbGetPrimaryDict() {
  return metadataStore.get('primaryDict', 'ccdict');
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

export async function dbUpdateWord(word, interval = 0, hasFlashCard = false) {
  console.log(`Adding new word: ${word}`);
  knex('words')
    .insert({
      word,
      interval,
      has_flash_card: hasFlashCard,
    })
    .onConflict('word')
    .merge(['interval', 'has_flash_card', 'updated_at'])
    .catch((err) => { console.log(err); });
}

// Insert words in chunks of chunkSize
export async function dbUpdateWords(wordRows) {
  try {
    await knex.transaction(async (trx) => {
      const chunkSize = 50;
      for (let i = 0; i < wordRows.length; i += chunkSize) {
        const chunk = wordRows.slice(i, i + chunkSize);
        await knex('words')
          .insert(chunk)
          .onConflict('word')
          .merge(['interval', 'has_flash_card', 'updated_at'])
          .transacting(trx);
      }
    });
  } catch (error) {
    console.log(error);
  }
}

export async function dbLoadWords() {
  const rows = await knex('words')
    .select({ word: 'word' })
    .catch((error) => { console.log(error); });
  const words = new Set();
  rows.forEach((row) => {
    words.add(row.word);
  });
  return words;
}

export async function dbWordExists(word) {
  const exists = await knex('words').select().where('word', word);
  return exists.length !== 0;
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
 *  filepath: string, // path of where book txt file is stored
 *  cover: string, // path of where book cover image is stored
 *  bookId: incrementing int,
 * }
 *
 * for each book there are also entries in the frequency table of their
 * word frequencies
 *
 ********************************** */
export async function dbAddBook(author, title, cover, filepath) {
  // For now just point to the actual txt file location in calibre.
  // Later we will make our own copy
  return knex('books').insert({
    author,
    title,
    cover,
    filepath,
  })
    .then(() => true)
    .catch((err) => {
      console.log(err);
      return false;
    });
}

export async function dbSaveWordTable(book, wordTable) {
  const wordRows = Object.entries(wordTable)
    .map(([word, frequency]) => ({
      book: book.bookId,
      word,
      count: frequency,
    }));
  // There should not be conflicts here.
  knex.batchInsert('frequency', wordRows, 100).catch((err) => {
    console.log(err);
  });
}

export async function dbLoadWordTable(book) {
  const wordRows = await knex('frequency')
    .select({ word: 'word', count: 'count' })
    .where('book', book.bookId);
  const wordDict = {};
  wordRows.forEach(({ word, count }) => {
    wordDict[word] = count;
  });
  return wordDict;
}

export async function dbBookSetCache(bookId, filepath) {
  return knex('books').where('bookId', bookId).update({
    segmented_file: filepath,
  });
}

const bookFields = {
  author: 'author',
  title: 'title',
  cover: 'cover',
  filepath: 'filepath',
  bookId: 'bookId',
  favorite: 'favorite',
  segmentedFile: 'segmented_file',
};

// Seems a bit repetative ...
export async function dbGetBooks(bookIds = []) {
  const books = knex('books').select(bookFields);
  if (bookIds.length > 0) {
    books.whereIn('bookId', bookIds);
  }
  return books;
}

export async function dbGetBook(author, title) {
  const books = await knex('books').select(
    bookFields,
  ).where({
    author, title,
  });
  return books[0];
}

export async function dbGetBookById(bookId) {
  const books = await knex('books').select(
    bookFields,
  ).where({
    bookId,
  });
  return books[0];
}

// For now we will use author and title to do book uniqueness
export async function dbBookExists(author, title) {
  const books = await knex('books').select(
    bookFields,
  ).where({
    author, title,
  });
  return books.length === 1;
}

export const databaseIpc = {
  getOptionValue, setOptionValue,
};
