import jieba from 'nodejieba';
import path from 'path';
// import { once } from 'events';
import fs from 'fs';
// import readline from 'readline';
import type {
  Book,
  SegmentedSentence,
} from './types';
// import { isInDictionary } from './dictionaries';
// direct from db to prevent cyclic dependency
import { dbGetBooks, dbBookSetCache } from './database';

const cache: {
  [path:string]: any
} = {};

const jiebaFiles = {
  dict: path.join(
    __dirname,
    '../../node_modules/nodejieba/dict/jieba.dict.utf8',
  ),
  hmmDict: path.join(
    __dirname,
    '../../node_modules/nodejieba/dict/hmm_model.utf8',
  ),
  userDict: path.join(
    __dirname,
    '../../node_modules/nodejieba/dict/user.dict.utf8',
  ),
  idfDict: path.join(
    __dirname,
    '../../node_modules/nodejieba/dict/idf.utf8',
  ),
  stopWordDict: path.join(
    __dirname,
    '../../node_modules/nodejieba/dict/stop_words.utf8',
  ),
};
async function computeDict(userConfigDir:string) {
  // TODO PORTOVER
  // Load a copy of the jieba dict
  // const prodDictFolder = path.join(process.resourcesPath, 'dict');
  // const inputFile = import.meta.env.MODE === 'production'
  //   ? path.join(prodDictFolder, 'jieba.dict.utf8')
  //   : './node_modules/nodejieba/dict/jieba.dict.utf8';
  // const outputFile = path.join(
  // app.getPath('userData'), 'jieba.mod.dict.utf8');
  // const inputStream = fs.createReadStream(inputFile);
  // const outputStream = fs.createWriteStream(
  // outputFile, { encoding: 'utf8' });
  // const lineReader = readline.createInterface({
  //   input: inputStream,
  //   terminal: false,
  // });
  // lineReader.on('line', (line) => {
  //   const items = line.split(' ');
  //   const [word] = items;
  //   if (isInDictionary(word)) {
  //     outputStream.write(`${line}\n`);
  //   }
  // });
  // await once(lineReader, 'close');
  //
  const copiedFiles = Object.fromEntries(
    Object.entries(jiebaFiles).map(
      ([key, filepath]) => {
        const filename = path.basename(filepath);
        const configDir = path.join(userConfigDir, '/Dictionaries');
        const newFile = path.join(configDir, filename);
        const fileContents = fs.readFileSync(filepath);
        fs.writeFileSync(newFile, fileContents);
        return [key, newFile];
      },
    ),
  );
  jieba.load(copiedFiles);
  // } else {
  //   jieba.load({
  //     dict: outputFile,
  //   });
  // }
}
export async function loadSegmentedText(book:Book) {
  if (!book.segmentedFile) {
    console.error('big problem');
  }
  return cache[book.segmentedFile];
}

async function doFullSegmentation(book:Book, userConfigDir:string) {
  // If the book has already been reduced to sentences previously
  if (book.segmentedFile) {
    const cacheLocation = path.join(
      userConfigDir,
      'segmentationCache',
      book.segmentedFile,
    );
    const sentenceSeg = await fs.promises.readFile(cacheLocation, {
      encoding: 'utf-8',
      flag: 'r',
    });
    const parsed = JSON.parse(sentenceSeg);
    cache[book.segmentedFile] = parsed;
  }
  // Otherwise we need to calculate the text
  const fullSegmentation = await loadJieba(book);
  console.log(typeof fullSegmentation);
  const joinedSentences = fullSegmentation.map(
    (sentence) => sentence.map(([word]) => word).join(''),
  );
  const fileName = `${book.title}-${book.author}.json`;
  const cacheLocation = path.join(
    userConfigDir,
    'segmentationCache',
    fileName,
  );
  await fs.promises.writeFile(cacheLocation, JSON.stringify(joinedSentences));
  dbBookSetCache(book.bookId, fileName);
}

export function segmentSentence(sentence:string):SegmentedSentence {
  const json = jieba.cut(sentence);

  return json.map((word) => {
    // const punc = /\p{Script_Extensions=Han}/u;
    // const punc = /\p{CJK_Symbols_and_Punctuation}/u;
    const punc = /\p{P}/u;
    if (punc.test(word)) {
      // punctuation
      return [word, 1];
    }
    if (/\s+/.test(word)) {
      // whitespace
      return [word, 1];
    }
    if (/\p{Script=Latin}+/u.test(word)) {
      // english
      return [word, 1];
    }
    if (/\p{Script=Han}+/u.test(word)) {
      return [word, 3];
    }
    // console.log(`unknown ${word}`);
    return [word, 1];
  });
}

export async function loadJieba(book:Book) {
  const txtPath = book.filepath;
  // console.log(`Loading ${txtPath} for the first time`);
  const txt = await fs.promises.readFile(txtPath, {
    encoding: 'utf-8',
    flag: 'r',
  });
  // Misses names, but also makes less compound words
  // Haha, I see why they recommended the default. This still produces a
  // 'lower' accuracy than CTA, but it is not as bad as others
  // const json = rsjieba.cut(txt);
  //
  const json = jieba.cut(txt);

  // Detects names better but makes stuff like 有庆死, 看凤霞
  // const json = nodejieba.cut(txt, true);

  // Creates weird words like 看家珍, 他们说
  // const json = nodejieba.cutHMM(txt);

  // Creates words like 两条腿
  // const json = nodejieba.cutAll(txt);

  // Doesn't get as many names still makes 两条腿
  // const json = nodejieba.cutForSearch(txt);

  return json.reduce(
    (result, origword) => {
      let type:number;
      let word = origword;
      // const punc = /\p{Script_Extensions=Han}/u;
      // const punc = /\p{CJK_Symbols_and_Punctuation}/u;
      const punc = /\p{P}/u;
      if (punc.test(word)) {
      // punctuation
        type = 1;
      } else if (/\s+/.test(word)) {
      // whitespace
        type = 1;
      } else if (/\p{Script=Latin}+/u.test(word)) {
      // english
        type = 1;
      } else if (/\p{Script=Han}+/u.test(word)) {
        type = 3;
      } else {
        type = 1;
      }
      const end = result[result.length - 1];
      if (word.length > 1 && word.includes('.')) {
      // It sees 15. 14. etc as being words,
      // so remove the . since it breaks db storage
        word = word.replaceAll('.', '');
      }
      if (word === '\n') {
        if (end.length > 0) {
          result.push([]);
        }
      } else if (
        word === '？'
        || word === '！'
        || word === '。'
        || word === '…'
        || word === '.'
      ) {
        if (end.length === 0) {
          const previous = result[result.length - 2];
          previous.push([word, type]);
        } else {
          end.push([word, type]);
          result.push([]);
        }
      } else if (word === ' ' || word === '　' || word === '\t') {
      // cta strips leading spaces
        if (end.length > 0) {
          end.push([word, type]);
        }
      } else if (
        (word === '”' || word === '‘' || word === '』')
        && end.length === 0
      ) {
      // Closing quotes go onto previous
        const previous = result[result.length - 2];
        previous.push([word, type]);
      } else {
        end.push([word, type]);
      }
      return result;
    },
    ([[]] as [string, number][][]),
  );
}

export async function preloadWords(userConfigDir:string) {
  await computeDict(userConfigDir);
  const books = await dbGetBooks();
  await Promise.all(books.map(
    (bookInfo) => doFullSegmentation(bookInfo, userConfigDir),
  ));
}
/* const TYPE = {
  NONE: 0, // None - Indicative of an error
  INVALID: 1, // Invalid - Invalid utf8 text
  CHINESE: 2, // Chinese - A word made up of Chinese text
  ALPHA: 3, // Alpha - A word made up of letters from the English alphabet
  NUMBER: 4, // Number - A word made up of Arabic numerals
  WHITESPACE: 5, // Whitespace - A block of whitespace
                 // (spaces, tabs, newlines etc).
  CHINESEPUNCTUATION: 6, // ChinesePunctuation - Chinese punctuation
  ASCIIPUNCTUATION: 7, // AsciiPunctuation - Standard ascii
                       // (English) punctuation.
}; */
