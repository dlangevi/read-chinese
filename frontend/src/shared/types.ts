import type { InjectionKey } from 'vue';

export type DictionaryInfo = {
  path: string,
  type: DictionaryType,
};

export type DictionaryType = 'english' | 'chinese';

export type DictionaryEntry = {
  definition: string;
  pronunciation: string;
};

export type KnownWords = {
  [key:string]: {
    interval:number;
  };
};

export type UnknownWordEntry = {
  word: string,
  occurance?: number,
  definition?: string,
  pinyin?: string,
};

export type HskVersion = '2.0' | '3.0';
export type HskLevel = 1 | 2 | 3 | 4 | 5 | 6 | 7;

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
};

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
};

export type UserSetting = {
  value:string;
  label:string;
  tooltip?:string;
  defaultValue:any;
  type:any
  other?:any;
  read?: any;
  write?:any;
  readFromBackEnd?:any;
  loaded?:boolean;
  cached?:any;
};

export type UserSettingsType = {
  [section:string]: {
    [label:string]: UserSetting;
  }
};
export const UserSettingsKey = Symbol('u') as InjectionKey<UserSettingsType>;
