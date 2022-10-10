export type dictionaryType = 'english' | 'chinese';

export type KnownWords = {
  [key:string]: {
    interval:number;
  };
};

export type SegmentedSentence = [token:string, type:number][];

export type BookStats = {
  probablyKnownWords: number;
  knownCharacters: number;
  totalCharacters: number;
  totalWords: number;
  totalKnownWords: number;
  targets: number[];
  targetOccurances: number[];
  needToKnow: number[];
}

export function initBookStats():BookStats {
  return {
    probablyKnownWords: 0,
    knownCharacters: 0,
    totalCharacters: 0,
    totalWords: 0,
    totalKnownWords: 0,
    targets: [],
    targetOccurances: [],
    needToKnow: [],
  };
}

export type Book = {
  author: string;
  title: string;
  cover: string;
  filepath: string;
  bookId: number;
  favorite: boolean;
  segmentedFile: string;
  hasRead: boolean;
  stats: BookStats;
}
