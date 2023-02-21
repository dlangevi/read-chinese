// TODO backend does not export these types
// even though IMO it should
import { backend } from '@wailsjs/models';

/*
 * type DictionaryInfoMap map[string]DictionaryInfo
 **/
export type DictionaryInfoMap = {
  [name:string] : backend.DictionaryInfo
};

/*
 * type WordDefinitions map[string]DictionaryDefinition
 **/
export type WordDefinitions = {
  [word:string] : backend.DictionaryDefinition
}
