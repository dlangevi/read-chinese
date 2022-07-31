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

export function generateSentences(words, modifyCards = false) {
  const sentences = {};
  const wordDict = {};
  const shouldLearn = {

  };
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

        candidates.push(toText(sentence));
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
  Object.entries(wordDict).forEach(([word, candidate]) => {
    if (candidate !== '') {
      goodOnes += 1;
      // For now do 10 at a time with lots of debugging
      if (modifyCards && goodOnes < 50) {
        addSentenceToCard(word, candidate);
      }
    }
  });
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

export function otherFunction() {
  console.log(';why');
}
