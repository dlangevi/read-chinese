// All saved data will be ran through here, so we can swap electron-store to
// something more performant later if we need
import Store from 'electron-store';

const wordStore = new Store({ name: 'words' });
const bookStore = new Store({ name: 'books' });
const metadataStore = new Store({ name: 'metadata' });

export function updateTimesRan() {
  const timesRan = metadataStore.get('ran', 0);
  metadataStore.set('ran', timesRan + 1);
}

export function getTimesRan() {
  return metadataStore.get('ran', 0);
}

export function saveWords(words) {
  wordStore.set('wordlist', words);
}

export function loadWords() {
  return wordStore.get('wordlist', {});
}

export function addBook(author, title, cover, filepath) {
  // For now just point to the actual txt file location in calibre. Later we will make our own copy
  console.log('Adding ', author, title, filepath);
  const books = bookStore.get('booklist', {});
  books[`${author}-${title}`] = {
    author,
    title,
    txtFile: filepath,
    cover,
  };
  bookStore.set('booklist', books);
}

export function getBooks() {
  return Object.values(bookStore.get('booklist', {}));
}

export function getBookKey(bookKey) {
  return bookStore.get('booklist')[bookKey];
}

// For now we will use author and title to do book uniqueness
export function bookExists(author, title) {
  const books = bookStore.get('booklist', {});
  return (`${author}-${title}` in books);
}
