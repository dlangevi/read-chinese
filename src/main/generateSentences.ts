import { getBooks } from './bookLibrary';
import { loadSegmentedText, segmentSentence } from './segmentation';
import { isKnown, updateInterval } from './knownWords';
import { SegmentedSentence } from '../shared/types';

function isT1Candidate(sentence:SegmentedSentence, t1word:string) {
  // TODO its possible the t1word is actually split across two neighbours, and
  // is not actually in the sentence. On the other hand, this can help in
  // finding extra instances of _some_ words
  return sentence.every(([word, type]) => {
    if (type !== 3) return true;
    if (word === t1word) return true;
    return isKnown(word);
  });
}

async function getSentencesForWord(word:string, {
  bookIds,
  skipBook,
}: { bookIds?: number[], skipBook?: number } = {}) {
  updateInterval();
  const books = await getBooks(bookIds);
  const candidates = new Set<string>();
  await Promise.all(books.map(async (bookInfo) => {
    if (bookInfo.bookId === skipBook) {
      return;
    }
    const fullSegmented = await loadSegmentedText(bookInfo);
    fullSegmented.forEach((sentence:string) => {
      if (sentence.includes(word)) {
        const segmented = segmentSentence(sentence);
        const justWords = segmented.map(([w]) => w);
        if (justWords.includes(word) && isT1Candidate(segmented, word)) {
          candidates.add(sentence);
        }
      }
    });
  }));
  const sentences = [...candidates];
  sentences.sort((a:string, b:string) => {
    // Ideal sentence is 20 characters, no real reason
    const aScore = Math.abs(a.length - 20);
    const bScore = Math.abs(b.length - 20);
    return aScore - bScore;
  });
  console.log(`generated ${sentences.length} sentences`);
  sentences.splice(10);
  return sentences;
}

export const generateSentencesIpc = { getSentencesForWord };
