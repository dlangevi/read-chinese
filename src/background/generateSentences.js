import { getBooks } from './bookLibrary';
import { loadJieba } from './segmentation';
import { isKnown } from './knownWords';

function toText(sentence) {
  return sentence.map(([word]) => word).join('');
}

function isT1Candidate(sentence, t1word) {
  return sentence.every(([word, type]) => {
    if (type !== 3) return true;
    if (word === t1word) return true;
    return isKnown(word);
  });
}

function sentenceKnown(sentence) {
  const unknowns = new Set();
  sentence.forEach(([word, type]) => {
    if (type !== 3) return;

    if (!(isKnown(word))) {
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
export async function whatShouldILearn(books = []) {
  if (books.length === 0) {
    books = await getBooks();
  }
  const shouldLearn = {};
  await Promise.all(books.map(async (bookInfo) => {
    const segmented = await loadJieba(bookInfo.filepath);
    segmented.forEach((sentence) => {
      const word = getT1Word(sentence);
      if (word) {
        if (!(word in shouldLearn)) {
          shouldLearn[word] = 0;
        }
        shouldLearn[word] += 1;
      }
    });
  }));
  const sorted = Object.entries(shouldLearn)
    .filter(([_, timesSeen]) => (timesSeen > 50))
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

async function getSentencesForWord(word, books = []) {
  if (books.length === 0) {
    books = await getBooks();
  }
  const candidates = new Set();
  await Promise.all(books.map(async (bookInfo) => {
    const segmented = await loadJieba(bookInfo.filepath);
    segmented.forEach((sentence) => {
      const text = toText(sentence);
      if (text.includes(word)) {
        if (isT1Candidate(sentence, word)) {
          candidates.add(text);
        }
      }
    });
  }));
  const sentences = [...candidates];
  sentences.sort((a, b) => (b.length - a.length));
  sentences.splice(10);
  return sentences;
}

// Generates sentences for a given set of words
// Returns a map of words to a sentence for each one
export async function generateSentences(
  words = [],
  books = [],
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
  await Promise.all(books.map(async (bookInfo) => {
    const segmented = await loadJieba(bookInfo.filepath);
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
  }));

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

export const generateSentencesIpc = { getSentencesForWord };
