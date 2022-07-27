import { getBooks } from './database';
import { loadJieba } from './segmentation';
import known from './knownWords';

function toText(sentence) {
  return sentence.map(([word]) => word).join('');
}

function sentenceKnown(sentence) {
  let unknown = 0;
  sentence.forEach(([word, type]) => {
    if (type !== 3) return;

    if (!(known.isKnown(word))) {
      unknown += 1;
    }
  });
  return unknown <= 0;
}

export function generateSentences(words) {
  const sentences = {};
  const wordDict = {};
  words.forEach((word) => {
    sentences[word] = [];
    wordDict[word] = '';
  });
  const allBooks = getBooks();
  const candidates = [];
  allBooks.forEach((bookInfo) => {
    console.log(`Loading ${bookInfo.txtFile}`);
    const segmented = loadJieba(bookInfo.txtFile);
    segmented.forEach((sentence) => {
      if (sentenceKnown(sentence)) {
        const text = toText(sentence);
        // For now longest wins
        sentence.forEach(([word, type]) => {
          if (type !== 3) return;

          const previousText = wordDict[word];
          if (previousText !== undefined && previousText.length < text.length) {
            wordDict[word] = text;
          }
        });

        candidates.push(toText(sentence));
      }
    });
  });
  let goodOnes = 0;
  Object.values(wordDict).forEach((candidate) => {
    if (candidate !== '') {
      goodOnes += 1;
    }
  });
  console.log(wordDict);
  console.log(`Generated ${goodOnes}/${Object.keys(wordDict).length}`);
}

export function otherFunction() {
  console.log(';why');
}
