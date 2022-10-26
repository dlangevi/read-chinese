import { Calibre } from 'node-calibre';
import { addBook, bookExists } from './bookLibrary';

interface CalibreBook {
  authors: string,
  cover: string,
  title: string,
  formats: [string],
}

// async function getCalibreBooks(calibreDir:string):Promise<CalibreBook[]> {
async function getCalibreBooks():Promise<CalibreBook[]> {
  // Create Calibre instance
  const calibre = new Calibre();
  const result = await calibre.run('calibredb list', [], {
    // limit: 10,
    forMachine: '',
    fields: 'cover,authors,title,formats',
  });
  const books : CalibreBook[] = JSON.parse(result);
  return books;
}

export async function importCalibreBooks() {
  // PORTOVER
  // const calibreDir = dialog.showOpenDialogSync({
  //   properties: ['openDirectory'],
  // });
  // if (calibreDir === undefined) {
  //   console.error('directory was not selected');
  //   return;
  // }
  // Technically this will work without selection since it
  // just does the default library last picked by calibre
  const books:CalibreBook[] = await getCalibreBooks();
  Object.values(books).forEach(async (book) => {
    // TODO do some conversion of epubs etc
    // (is epub-convert avaliable on all platforms?)
    if (!(await bookExists(book.authors, book.title))) {
      console.log(`Creating book ${book.authors} ${book.title}`);
      const txtBooks:string[] = book.formats.filter(
        (path:string) => path.match(/.*.txt/),
      );
      if (txtBooks.length > 0) {
        addBook(book.authors, book.title, book.cover, txtBooks[0]);
      }
    }
  });
  return 'finish';
}

export const calibreIpc = {
  importCalibreBooks,
};
