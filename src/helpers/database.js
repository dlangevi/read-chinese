// All saved data will be ran through here, so we can swap electron-store to
// something more performant later if we need
import Store from 'electron-store';
import Knex from 'knex';
import knexConfigMap from '../../knexfile';

const knexConfig = knexConfigMap[process.env.NODE_ENV];
const knex = Knex(knexConfig);

knex.migrate.latest(knexConfig).catch((err) => {
  console.log(err);
});

function bookKey(author, title) {
  return `${author}-${title}`;
}

// TODO should this also be in sql?
const bookStore = new Store({ name: 'books' });
const metadataStore = new Store({ name: 'metadata' });

export function updateTimesRan() {
  const timesRan = metadataStore.get('ran', 0);
  metadataStore.set('ran', timesRan + 1);
}

export function getTimesRan() {
  return metadataStore.get('ran', 0);
}

export async function updateWord(ankiCard) {
  const word = ankiCard.fields.Hanzi.value;
  const exists = await knex('words').select().where('word', word);
  if (exists.length === 0) {
    console.log(`Adding new word: ${word}`);
    knex('words')
      .insert({
        word,
        interval: ankiCard.interval,
        has_flash_card: true,
      }).catch((err) => { console.log(err); });
  } else {
    knex('words')
      .where('word', word)
      .update({
        has_flash_card: true,
        interval: ankiCard.interval,
      }).catch((err) => { console.log(err); });
  }
}

export function saveWords(words) {
  Object.entries(words).forEach(([word, entry]) => {
    console.log(`Inserting ${word}`);
    knex('words')
      .insert({
        word,
        interval: entry.interval,
      })
      .catch((err) => {
        console.error(err);
      });
  });
}

export async function wordExists(word) {
  const exists = await knex('words')
    .select().where('word', word)
    .catch((err) => {
      console.error(err);
    });
  return exists.length !== 0;
}

// TODO more arguments
export async function saveWord(word, interval = 0) {
  const exists = await wordExists(word);
  if (!exists) {
    knex('words')
      .insert({
        word,
        interval,
      })
      .catch((err) => {
        console.error(err);
      });
  }
}

export async function loadWords() {
  const rows = await knex('words')
    .select({ id: 'id', word: 'word', interval: 'interval' })
    .catch((error) => { console.log(error); });
  const words = {};
  rows.forEach((row) => {
    words[row.word] = {
      interval: row.interval,
      // todo get added working
      added: 0,
    };
  });
  return words;
}

export function addBook(author, title, cover, filepath) {
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

export function getBooks() {
  return Object.values(bookStore.get('booklist', {}));
}

export function getBook(author, title) {
  return bookStore.get('booklist')[bookKey(author, title)];
}

export function getBookByID(bookID) {
  return bookStore.get('booklist')[bookID];
}

// For now we will use author and title to do book uniqueness
export function bookExists(author, title) {
  const books = bookStore.get('booklist', {});
  return (bookKey(author, title) in books);
}
