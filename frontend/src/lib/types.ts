import type { InjectionKey } from 'vue';

export type UserSetting = {
  name:string;
  label:string;
  tooltip?:string;
  type:any
  read?: any;
  write?:any;
};

export type UserSettingsType = {
  [section:string]: {
    [label:string]: UserSetting;
  }
};
export const UserSettingsKey = Symbol('u') as InjectionKey<UserSettingsType>;

/// //////////////////////////////////////////////// ///
//  Shared Types! Keep in sync with src-node/types.ts //
/// //////////////////////////////////////////////// ///
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
