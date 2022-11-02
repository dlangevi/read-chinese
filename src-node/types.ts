export type KnownWords = {
  [key:string]: {
    interval:number;
  };
};
export type SegmentedSentence = [token:string, type:number][];

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

/// //////////////////////////////////////////////////////// ///
//  Shared Types! Keep in sync with frontend/src/lib/types.ts //
/// //////////////////////////////////////////////////////// ///
export type DictionaryInfo = {
  path: string,
  type: DictionaryType,
};

export type DictionaryType = 'english' | 'chinese';

export type DictionaryEntry = {
  definition: string;
  pronunciation: string;
};

export type UnknownWordEntry = {
  word: string,
  occurance?: number,
  definition?: string,
  pinyin?: string,
};

export type HskVersion = '2.0' | '3.0';
export type HskLevel = 1 | 2 | 3 | 4 | 5 | 6 | 7;

export type BookStats = {
  probablyKnownWords: number;
  knownCharacters: number;
  totalCharacters: number;
  totalWords: number;
  totalKnownWords: number;
  targets: number[];
  targetOccurances: number[];
  needToKnow: number[];
};

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
};
