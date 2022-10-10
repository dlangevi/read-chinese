// Dont know if this is a silly idea, but want to isolate all the actual reads
// and writes in a seperate file. So anything that is saved or read from
// persistant storage will have to go through functions prefixed 'db' in this
// file. If we want to swap out backends later on at least all the code to be
// changed will be in the same place
import Store, { Schema } from 'electron-store';
import Knex from 'knex';
// For now we do the sync whenever the db changes.
import knexConfigMap from '../../knexfile.mjs';
import {
  dictionaryType, KnownWords, Book,
} from '../shared/types';

console.log(knexConfigMap);
const knexConfig = knexConfigMap[import.meta.env.MODE];
export const knex = Knex(knexConfig);

// This is called and awaited before before anyother code can run
export async function initializeDatabase() {
  await knex.migrate.latest().catch((err) => {
    console.log(err);
  });
}

interface metadata {
  ran: number;
  dicts: {
    [name:string] : {
    path: string,
    type: dictionaryType,
    }
  };
  primaryDict: string
}
const metadataSchema : Schema<metadata> = {
  ran: {
    type: 'number',
    default: 0,
  },
  dicts: {
    type: 'object',
    default: {},
  },
  primaryDict: {
    type: 'string',
    default: 'ccdict',
  },
};

// Dictionaries and User settings can be stored in electron-store
// since they are low footprint
const metadataStore = new Store({
  name: 'metadata',
  schema: metadataSchema,
});

/** *********************************
 *
 * Metadata
 *
 ********************************** */

export function updateTimesRan():void {
  const timesRan:number = metadataStore.get('ran');
  metadataStore.set('ran', timesRan + 1);
}

export function getTimesRan():number {
  return metadataStore.get('ran');
}

// TODO, should there be some white list of valid keys?
// import schema type info from usersettings.ts, including
// default values
export function getOptionValue(key:string, defaultValue:any) {
  return metadataStore.get(key, defaultValue);
}

function setOptionValue(key:string, value:any) {
  metadataStore.set(key, value);
}

export function dbSaveDict(name:string, path:string, type:dictionaryType) {
  const dicts = metadataStore.get('dicts');
  dicts[name] = {
    path,
    type,
  };
  metadataStore.set('dicts', dicts);
}

export function dbDeleteDict(name:string) {
  const dicts = metadataStore.get('dicts');
  delete dicts[name];
  metadataStore.set('dicts', dicts);
}

export function dbLoadDicts() {
  return metadataStore.get('dicts');
}

export function dbSetPrimaryDict(dictName:string) {
  return metadataStore.set('primaryDict', dictName);
}
export function dbGetPrimaryDict() {
  return metadataStore.get('primaryDict');
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

export async function dbUpdateWord(
  word:string,
  interval = 0,
  hasFlashCard = false,
) {
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
export async function dbUpdateWords(wordRows:any) {
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

interface Word {
  word: string;
  has_flash_card: boolean;
  has_sentence: boolean;
  interval: number;
  created_at: Date;
  updated_at: Date;
}

export async function dbLoadWords() {
  return knex<Word>('words')
    .select({ word: 'word', interval: 'interval' })
    .then((rows) => {
      const words: KnownWords = {};
      rows.forEach((row) => {
        words[row.word] = { interval: row.interval };
      });
      return words;
    })
    .catch((error) => {
      console.log(error);
      return error;
    });
}

export async function dbWordExists(word:string) {
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
export async function dbAddBook(
  author:string,
  title:string,
  cover:string,
  filepath:string,
) {
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

export async function dbSaveWordTable(book:any, wordTable:any) {
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

export async function dbLoadWordTable(book:any) {
  const wordRows = await knex('frequency')
    .select({ word: 'word', count: 'count' })
    .where('book', book.bookId);
  const wordDict: {[key:string]: number } = {};
  wordRows.forEach(({ word, count }) => {
    wordDict[word] = count;
  });
  return wordDict;
}

export async function dbBookSetCache(bookId:number, filepath:string) {
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
  hasRead: 'has_read',
};

// Seems a bit repetative ...
export async function dbGetBooks(bookIds:number[] = []) :Promise<Book[]> {
  const books = knex('books').select(bookFields);
  if (bookIds.length > 0) {
    books.whereIn('bookId', bookIds);
  }
  return books;
}

export async function dbGetBook(author:string, title:string):Promise<Book> {
  const books = await knex('books').select(
    bookFields,
  ).where({
    author, title,
  });
  return books[0];
}

export async function dbGetBookById(bookId:number):Promise<Book> {
  const books = await knex('books').select(
    bookFields,
  ).where({
    bookId,
  });
  return books[0];
}

// For now we will use author and title to do book uniqueness
export async function dbBookExists(
  author:string,
  title:string,
):Promise<Boolean> {
  const books = await knex('books').select(
    bookFields,
  ).where({
    author, title,
  });
  if (books.length !== 1) {
    console.log(books);
  }
  return books.length === 1;
}

export const databaseIpc = {
  getOptionValue, setOptionValue,
};
