import { inject } from 'vue';
import type { InjectionKey } from 'vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import SettingsTextbox
  from '@/components/SettingsWidgets/SettingsTextbox.vue';
import DictionariesList
  from '@/components/SettingsWidgets/DictionariesList.vue';
import ModelManager
  from '@/components/SettingsWidgets/ModelManager.vue';
import SettingsSlider
  from '@/components/SettingsWidgets/SettingsSlider.vue';
import {
  GetUserSetting,
  SetUserSetting,
  GetUserSettingBool,
  SetUserSettingBool,
  GetUserSettingInt,
  SetUserSettingInt,
} from '@wailsjs/backend/UserSettings';

type UserSetting = {
  name:string;
  label:string;
  tooltip?:string;
  type:any
  read?: any;
  write?:any;
};

async function settingsObject(
  name:string,
  label:string,
  tooltip:string,
  widgetType:any,
  getter:any,
  setter:any,
):Promise<UserSetting> {
  let value = await getter(name);
  const option:UserSetting = {
    name,
    label,
    tooltip,
    type: widgetType,
  };

  option.read = function read() {
    return value;
  };
  option.write = async function write(newValue:any) {
    value = newValue;
    return setter(name, value);
  };
  return option;
}

async function checkBox(
  name:string,
  label:string,
  tooltip:string,
) {
  return settingsObject(
    name,
    label,
    tooltip,
    SettingsCheckbox,
    GetUserSettingBool,
    SetUserSettingBool,
  );
}

async function textBox(
  name:string,
  label:string,
  tooltip:string,
) {
  return settingsObject(
    name,
    label,
    tooltip,
    SettingsTextbox,
    GetUserSetting,
    SetUserSetting,
  );
}

async function slider(
  name:string,
  label:string,
  tooltip:string,
) {
  return settingsObject(
    name,
    label,
    tooltip,
    SettingsSlider,
    GetUserSettingInt,
    SetUserSettingInt,

  );
}

async function loadSettings(settings : Promise<UserSetting>[]) {
  const waited = await Promise.all(settings);
  const options :{ [label:string]:UserSetting } = {};
  waited.forEach((setting) => {
    options[setting.name] = setting;
  });
  return options;
}

// export async function generateUserSettings() :Promise<UserSettingsType> {
export async function generateUserSettings() {
  const CardCreation = await loadSettings([
    checkBox(
      'AutoAdvanceSentence',
      'Auto advance after sentence selection',
      'After picking a sentence, move to the next step',
    ),
    checkBox(
      'PopulateEnglish',
      'Auto fill english definitions',
      'If only one definition exists, auto select it',
    ),
    checkBox(
      'PopulateChinese',
      'Auto fill chinese definitions',
      'If only one definition exists, auto select it',
    ),
    checkBox(
      'AutoAdvanceEnglish',
      'Auto advance after definition selection',
      'After picking a definition, move to the next step',
    ),
    checkBox(
      'AutoAdvanceImage',
      'Auto advance after image selection',
      'After picking a image, move to the next step',
    ),
    checkBox(
      'AutoAdvanceCard',
      'Create card once all fields have been filled',
      'Create card once all fields have been filled',
    ),

  ]);
  const Dictionaries = await loadSettings([
    (async () => ({
      name: 'Dictionaries',
      label: 'Dictionaries',
      type: DictionariesList,
    } as UserSetting))(),
    checkBox(
      'ShowDefinitions',
      'Show Definitions',
      'Show the definitions for words in various tables',
    ),
    checkBox(
      'EnableChinese',
      'Use Chinese definitions',
      'Allow flashcards to use chinese ' +
      'definitions instead of just english ones',
    ),
  ]);

  const SentenceGeneration = await loadSettings([
    slider(
      'KnownInterval',
      'Time before a word is considered "known"',
      'How long of an interval in anki before a word is ' +
      ' included in generated sentences',
    ),
    slider(
      'IdealSentenceLength',
      'Ideal Sentence Length"',
      'What the ideal sentence length you want to be selected from books',
    ),

  ]);

  const AnkiConfig = await loadSettings([
    (async () => ({
      name: 'ModelManager',
      label: 'ModelManager',
      type: ModelManager,
    } as UserSetting))(),
    // TODO these should be read only? Need some better way of
    // seperating the Values from their associated display
    textBox(
      'ActiveDeck',
      'Active Deck',
      'The Deck we are loading and saving cards to',
    ),
    textBox(
      'ActiveModel',
      'Active Note Type',
      'The anki note type we are mapping the fields to',
    ),
    checkBox(
      'AddProgramTag',
      'Add read-chinese tag',
      'I',
    ),
    checkBox(
      'AddBookTag',
      'Add source book title tag',
      'Obck',
    ),
    checkBox(
      'AllowDuplicates',
      'Allow Duplicates',
      'K',
    ),
    checkBox(
      'GenerateTermAudio',
      'Auto generate audio for keyword',
      'Not implemented yet',
    ),
    checkBox(
      'GenerateSentenceAudio',
      'Auto generate audio for example sentence',
      'Not implemented yet',
    ),
    textBox(
      'AzureApiKey',
      'Azure Audio Api Key',
      'Setup an free azure tts account and put your key here',
    ),
    textBox(
      'AzureImageApiKey',
      'Azure Image Api Key',
      'Setup an free azure bing search and put your key here',
    ),
  ]);

  const BookLibrary = await loadSettings([
    checkBox(
      'OnlyFavorites',
      'Only show favorited books',
      'Obv',
    ),
    checkBox(
      'HideRead',
      'Hide read books',
      'Obv',
    ),
  ]);

  const CardManagement = await loadSettings([
    checkBox(
      'ProblemFlagged',
      'Flagged Cards',
      'exe',
    ),
    checkBox(
      'ProblemMissingImage',
      'Missing Images',
      'exe',
    ),
    checkBox(
      'ProblemMissingSentence',
      'Missing Sentence',
      'exe',
    ),
    checkBox(
      'ProblemMissingSentenceAudio',
      'Missing Sentence Audio',
      'eexexe',
    ),
    checkBox(
      'ProblemMissingWordAudio',
      'Missing Word Audio',
      'exe',
    ),
    checkBox(
      'ProblemMissingPinyin',
      'Missing Pinyin',
      'exe',
    ),
  ]);

  return {
    CardCreation,
    AnkiConfig,
    SentenceGeneration,
    Dictionaries,
    BookLibrary,
    CardManagement,
  };
}

type UserSettingsType = Awaited<ReturnType<typeof generateUserSettings>>;
export const UserSettingsKey = Symbol('u') as InjectionKey<UserSettingsType>;

export function getUserSettings():UserSettingsType {
  return inject(UserSettingsKey) as UserSettingsType;
}
