import jieba from 'nodejieba';
import path from 'path';
import { once } from 'events';
import fs from 'fs';
import readline from 'readline';
import { isInDictionary } from './dictionaries';
// direct from db to prevent cyclic dependency
import { dbGetBooks } from './database';

const cache = { };

async function computeDict() {
  // Load a copy of the jieba dict
  const inputStream = fs.createReadStream('./dict/jieba.dict.utf8');
  const outputStream = fs.createWriteStream(
    './dict/mod.dict.utf8',
    { encoding: 'utf8' },
  );
  const lineReader = readline.createInterface({
    input: inputStream,
    terminal: false,
  });
  lineReader.on('line', (line) => {
    const items = line.split(' ');
    const [word] = items;
    if (isInDictionary(word)) {
      outputStream.write(`${line}\n`);
    }
  });
  await once(lineReader, 'close');

  if (process.env.NODE_ENV === 'production') {
    // The default dict doesn't load from the asar archive for some reason
    // If in production use the copies we have made in resources
    jieba.load({
      dict: './dict/mod.dict.utf8',
      hmmDict: path.join(process.resourcesPath, './dict/hmm_model.utf8'),
      userDict: path.join(process.resourcesPath, './dict/user.dict.utf8'),
      idfDict: path.join(process.resourcesPath, './dict/idf.utf8'),
      stopWordDict: path.join(process.resourcesPath, './dict/stop_words.utf8'),
    });
  } else {
    jieba.load({
      dict: './dict/mod.dict.utf8',
    });
  }
}

export async function loadJieba(txtPath) {
  if (txtPath in cache) {
    return cache[txtPath];
  }
  // console.log(`Loading ${txtPath} for the first time`);
  const txt = await fs.promises.readFile(txtPath, 'UTF-8', 'r');
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

  const finalResult = json.reduce((result, origword) => {
    let type = '';
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
      // It sees 15. 14. etc as being words, so remove the . since it breaks db
      // storage
      word = word.replaceAll('.', '');
    }
    if (word === '\n') {
      if (end.length > 0) {
        result.push([]);
      }
    } else if (word === '？' || word === '！' || word === '。'
      || word === '…' || word === '.') {
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
    } else if ((word === '”' || word === '‘' || word === '』')
      && end.length === 0) {
      // Closing quotes go onto previous
      const previous = result[result.length - 2];
      previous.push([word, type]);
    } else {
      end.push([word, type]);
    }
    return result;
  }, [[]]);
  cache[txtPath] = finalResult;
  return finalResult;
}

export async function preloadWords() {
  await computeDict();
  const books = await dbGetBooks();
  books.forEach((bookInfo) => {
    loadJieba(bookInfo.filepath);
  });
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
