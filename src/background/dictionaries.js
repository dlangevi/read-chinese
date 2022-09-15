// load dictionary
import fs from 'fs';
import { dbSaveDictPath, dbLoadDictPaths } from './database';

const dicts = {};

export function addDictionary(path) {
  // Just for now
  // Later we will make our own copy of the dictionary
  dbSaveDictPath('ccdict', path);
}

export function loadDictionaries() {
  const dictPaths = dbLoadDictPaths();
  Object.entries(dictPaths).forEach(([name, path]) => {
    console.log(name, path);
    const fileContents = fs.readFileSync(path);
    const contents = JSON.parse(fileContents);
    const dictionary = {};
    contents.forEach((term) => {
      const word = term.term;
      if (!(word in dictionary)) {
        dictionary[word] = [];
      }
      dictionary[word].push(term);
    });
    dicts[name] = dictionary;
  });
}

export function getDefinition(word) {
  const term = dicts.ccdict[word];
  if (term === undefined) {
    return undefined;
  }
  return term[0].definition;
}
function getDefinitionsForWord(word) {
  const term = dicts.ccdict[word];
  if (term === undefined) {
    return [];
  }
  console.log(term);
  return term.map((def) => ({
    definition: def.definition,
    pronunciation: def.pronunciation,
  }));
}

export const dictionariesIpc = {
  getDefinitionsForWord,
};
