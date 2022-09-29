import { getBooks } from './bookLibrary';
import { loadSegmentedText, segmentSentence } from './segmentation';
import { isKnown } from './knownWords';

function isT1Candidate(sentence, t1word) {
  // TODO its possible the t1word is actually split across two neighbours, and
  // is not actually in the sentence. On the other hand, this can help in
  // finding extra instances of _some_ words
  return sentence.every(([word, type]) => {
    if (type !== 3) return true;
    if (word === t1word) return true;
    return isKnown(word);
  });
}

async function getSentencesForWord(word, {
  bookIds,
  skipBook,
} = {}) {
  const books = await getBooks(bookIds);
  const candidates = new Set();
  await Promise.all(books.map(async (bookInfo) => {
    if (bookInfo.bookId === skipBook) {
      return;
    }
    const fullSegmented = await loadSegmentedText(bookInfo);
    fullSegmented.forEach((sentence) => {
      if (sentence.includes(word)) {
        const segmented = segmentSentence(sentence);
        if (isT1Candidate(segmented, word)) {
          candidates.add(sentence);
        }
      }
    });
  }));
  const sentences = [...candidates];
  sentences.sort((a, b) => {
    // Ideal sentence is 20 characters, no real reason
    const aScore = Math.abs(a.length - 20);
    const bScore = Math.abs(b.length - 20);
    return aScore - bScore;
  });
  sentences.splice(10);
  return sentences;
}

export const generateSentencesIpc = { getSentencesForWord };
