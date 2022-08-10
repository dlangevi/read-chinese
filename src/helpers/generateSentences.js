import { getBooks } from './database';
import { loadJieba } from './segmentation';
import known from './knownWords';
import { addSentenceToCard } from './ankiInterface';

function toText(sentence) {
  return sentence.map(([word]) => word).join('');
}

function sentenceKnown(sentence) {
  const unknowns = new Set();
  sentence.forEach(([word, type]) => {
    if (type !== 3) return;

    if (!(known.isKnown(word))) {
      unknowns.add(word);
    }
  });
  return [...unknowns];
}

function getT1Word(sentence) {
  const unknowns = sentenceKnown(sentence);
  if (unknowns.length === 1) {
    return unknowns[0];
  }
  return false;
}

// Generates a list of words that are
// holding you back from T1 sentences
export function whatShouldILearn(books = []) {
  if (books.length === 0) {
    books = getBooks();
  }
  const shouldLearn = {};
  books.forEach((bookInfo) => {
    console.log(`Loading ${bookInfo.txtFile}`);
    const segmented = loadJieba(bookInfo.txtFile);
    segmented.forEach((sentence) => {
      const word = getT1Word(sentence);
      if (word) {
        if (!(word in shouldLearn)) {
          shouldLearn[word] = 0;
        }
        shouldLearn[word] += 1;
      }
    });
  });
  const sorted = Object.entries(shouldLearn)
    .filter(([_, timesSeen]) => (timesSeen > 100))
    .sort(([_, timesA], [__, timesB]) => {
      if (timesA > timesB) {
        return 1;
      }
      return 0;
    })
    .map(([word, timesSeen]) => ({
      word, occurance: timesSeen,
    }));
  return sorted;
}

// Generates sentences for a given set of words
// Returns a map of words to a sentence for each one
export async function generateSentences(
  words = [],
  books = [],
  modifyCards = false,
) {
  const wordDict = {};
  const shouldLearn = {

  };
  words.forEach((word) => {
    wordDict[word] = '';
  });
  if (books.length === 0) {
    books = getBooks();
  }
  books.forEach((bookInfo) => {
    console.log(`Loading ${bookInfo.txtFile}`);
    const segmented = loadJieba(bookInfo.txtFile);
    segmented.forEach((sentence) => {
      const unknowns = sentenceKnown(sentence);
      if (unknowns.length === 0) {
        const text = toText(sentence);
        // For now longest wins
        sentence.forEach(([word, type]) => {
          if (type !== 3) return;

          const previousText = wordDict[word];
          if (previousText !== undefined && previousText.length < text.length) {
            wordDict[word] = text;
          }
        });
      } else if (unknowns.length === 1) {
        const learn = unknowns[0];
        if (!(learn in shouldLearn)) {
          shouldLearn[learn] = 0;
        }
        shouldLearn[learn] += 1;
      }
    });
  });
  let goodOnes = 0;
  const entries = Object.entries(wordDict);

  for (let i = 0; i < entries.length; i += 1) {
    const [word, candidate] = entries[i];
    if (candidate !== '') {
      goodOnes += 1;
      // For now do 10 at a time with lots of debugging
      if (modifyCards) {
        // eslint-ignore
        await addSentenceToCard(word, candidate);
      }
    }
  }
  console.log(`Generated ${goodOnes}/${Object.keys(wordDict).length}`);
  const sorted = Object.entries(shouldLearn)
    .filter(([_, timesSeen]) => (timesSeen > 100))
    .sort(([_, timesA], [__, timesB]) => {
      if (timesA > timesB) {
        return 1;
      }
      return 0;
    });
  console.log(sorted);
}

export function initWordGenIpc(ipcMain) {
  ipcMain.handle('learningTarget', () => {
    const words = whatShouldILearn();
    return words;
  });
}
