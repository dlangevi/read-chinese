// load dictionary
import fs from 'fs';
import { dbSaveDict, dbLoadDicts } from './database';

const dicts = {
  english: {},
  chinese: {},
};
// TODO have these user set
const defaultType = 'english';
const defaultDict = 'ccdict';

export function addDictionary(name, path, type) {
  // Just for now
  // Later we will make our own copy of the dictionary
  dbSaveDict(name, path, type);
  // After adding a new dict load it
  loadDictionaries();
}

export function dictionaryInfo() {
  return dbLoadDicts();
}

export function loadDictionaries() {
  const ldicts = dbLoadDicts();
  Object.entries(ldicts).forEach(([name, entry]) => {
    console.log(name, entry);
    const fileContents = fs.readFileSync(entry.path);
    const contents = JSON.parse(fileContents);
    const dictionary = {};
    contents.forEach((term) => {
      const word = term.term;
      if (!(word in dictionary)) {
        dictionary[word] = [];
      }
      dictionary[word].push(term);
    });
    dicts[entry.type][name] = dictionary;
  });
}

// This is just used for a simple definition when displaying
// in large word lists
export function getDefaultDefinition(word) {
  const term = dicts[defaultType][defaultDict][word];
  if (term === undefined) {
    return undefined;
  }
  return term[0].definition;
}

export function getPinyin(word) {
  const terms = dicts[defaultType][defaultDict][word];
  if (!terms) {
    // TODO do char by char lookup and concatinate?
    return '';
  }
  if (terms.length === 0) {
    return terms[0].pronunciation;
  }
  return [...new Set(terms.map((term) => term.pronunciation))].join(', ');
}

export function isInDictionary(word) {
  return word in dicts[defaultType][defaultDict];
}

// type = 'english' or 'chinese'
function getDefinitionsForWord(word, type) {
  const filteredDicts = dicts[type];
  const answers = [];

  Object.values(filteredDicts).forEach((dict) => {
    const term = dict[word];
    if (term === undefined) {
      return;
    }
    term.forEach((def) => {
      answers.push({
        definition: def.definition.replace(/\n/g, '<br>'),
        pronunciation: def.pronunciation,
      });
    });
  });

  return answers;
}

export const dictionariesIpc = {
  getDefinitionsForWord,
  dictionaryInfo,
  addDictionary,
};
