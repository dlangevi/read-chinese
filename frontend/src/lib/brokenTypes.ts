// TODO backend does not export these types
// even though IMO it should
import { backend } from '@wailsjs/models';

/*
 * type DictionaryInfoMap map[string]DictionaryInfo
 * type DictionaryInfo struct {
 *   Name     string `json:"name"`
 *   Path     string `json:"path"`
 *   Language string `json:"type"`
 * }
 **/
export type DictionaryInfo = {
  name: string,
  path: string,
  type: string,
};
export type DictionaryInfoMap = {
  [name:string] : DictionaryInfo
};

/*
 * type WordDefinitions map[string]DictionaryDefinition
 *
 **/
export type WordDefinitions = {
  [word:string] : backend.DictionaryDefinition
}
