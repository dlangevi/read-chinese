// load dictionary
import fs from 'fs';
import {
  dbSaveDict, dbLoadDicts, dbGetPrimaryDict, dbSetPrimaryDict, dbDeleteDict,
} from './database';

const dicts = {
};
// TODO have these user set
let defaultDict = 'ccdict';

export function addDictionary(name, path, type) {
  // Just for now
  // Later we will make our own copy of the dictionary
  dbSaveDict(name, path, type);
  // After adding a new dict load it
  loadDictionaries();
}

export function deleteDictionary(name) {
  dbDeleteDict(name);
  delete dicts[name];
}

function setPrimaryDict(dictName) {
  console.log('setting primary dict to ', dictName);
  defaultDict = dictName;
  dbSetPrimaryDict(dictName);
}

export function dictionaryInfo() {
  return dbLoadDicts();
}

export function loadDictionaries() {
  const ldicts = dbLoadDicts();
  defaultDict = dbGetPrimaryDict();
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
    dicts[name] = {
      dictionary,
      type: entry.type,
    };
  });
}

// This is just used for a simple definition when displaying
// in large word lists
export function getDefaultDefinition(word) {
  const term = dicts[defaultDict].dictionary[word];
  if (term === undefined) {
    return undefined;
  }
  return term[0].definition;
}

export function getPinyin(word) {
  const terms = dicts[defaultDict].dictionary[word];
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
  return word in dicts[defaultDict].dictionary;
}

// type = 'english' or 'chinese'
function getDefinitionsForWord(word, type) {
  const answers = [];

  Object.values(dicts)
    .filter((dict) => dict.type === type)
    .forEach((dict) => {
      const term = dict.dictionary[word];
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
  setPrimaryDict,
  deleteDictionary,
};
