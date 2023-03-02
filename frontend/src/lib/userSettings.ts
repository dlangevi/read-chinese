import { inject, reactive } from 'vue';
import { LoadModels, LoadDecks } from '@wailsjs/backend/ankiInterface';
import type { InjectionKey } from 'vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import SettingsTextbox
  from '@/components/SettingsWidgets/SettingsTextbox.vue';
import SettingsSelector
  from '@/components/SettingsWidgets/SettingsSelector.vue';
import DictionariesList
  from '@/components/SettingsWidgets/DictionariesList.vue';
import ModelManager
  from '@/components/SettingsWidgets/ModelManager.vue';
import SettingsSlider
  from '@/components/SettingsWidgets/SettingsSlider.vue';
import type { backend } from '@wailsjs/models';

import {
  SetUserSetting,
  SetUserSettingBool,
  SetUserSettingInt,
  GetUserSettings,
} from '@wailsjs/backend/UserConfig';

export type UserSetting = {
  name:string;
  label:string;
  tooltip?:string;
  type:any;
  write:any;
  dataSource?: () => Promise<any>;
};

type UserConfigDisplay = {
  [K in keyof backend.UserConfig]: {
    [P in keyof backend.UserConfig[K]]?: UserSetting
  }
}

export let updateSettings = async () => {};
function settingsObject(
  name:string,
  label:string,
  tooltip:string|undefined,
  widgetType:any,
  setter:any,
):UserSetting {
  const option:UserSetting = {
    name,
    label,
    tooltip,
    type: widgetType,
    write: async function write(newValue:any) {
      console.log('writing', name, newValue);
      await setter(name, newValue);
      await updateSettings();
    },
  };
  return option;
}

function selector(
  name:string,
  label:string,
  dataSource: () => Promise<string[]>,
  tooltip?:string,
):UserSetting {
  const option:UserSetting = {
    name,
    label,
    tooltip,
    type: SettingsSelector,
    write: async function write(newValue:string) {
      console.log('writing', name, newValue);
      await SetUserSetting(name, newValue);
      await updateSettings();
    },
    dataSource,
  };
  return option;
}

function checkBox(
  name:string,
  label:string,
  tooltip?:string,
) {
  return settingsObject(
    name,
    label,
    tooltip,
    SettingsCheckbox,
    SetUserSettingBool,
  );
}

function textBox(
  name:string,
  label:string,
  tooltip?:string,
) {
  return settingsObject(
    name,
    label,
    tooltip,
    SettingsTextbox,
    SetUserSetting,
  );
}

function slider(
  name:string,
  label:string,
  tooltip?:string,
) {
  return settingsObject(
    name,
    label,
    tooltip,
    SettingsSlider,
    SetUserSettingInt,
  );
}

function loadSettings(settings : UserSetting[]) {
  const options :{ [label:string]:UserSetting } = {};
  settings.forEach((setting) => {
    if (setting.name !== undefined) {
      options[setting.name] = setting;
    }
  });
  return options;
}

export const ComponentTable = loadSettings([
  selector(
    'Theme',
    'Select Theme',
    async function () {
      return [
        'light', 'dark', 'cupcake',
        'bumblebee', 'emerald', 'corporate',
        'synthwave', 'retro', 'cyberpunk',
        'valentine', 'halloween', 'garden',
        'forest', 'aqua', 'lofi', 'pastel',
        'fantasy', 'wireframe', 'black',
        'luxury', 'dracula', 'cmyk',
        'autumn', 'business', 'acid',
        'lemonade', 'night', 'coffee', 'winter',
      ];
    },
  ),
  checkBox(
    'EnableExperimental',
    'Enable Experimental Features',
    'Only do this if you are me',
  ),
  selector(
    'ActiveDeck',
    'Active Anki Deck',
    LoadDecks,
  ),
  selector(
    'ActiveModel',
    'Active Anki Model',
    LoadModels,
  ),
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
    'Auto advance after english definition selection',
    'After picking an english definition, move to the next step',
  ),
  checkBox(
    'AutoAdvanceChinese',
    'Auto advance after chinese definition selection',
    'After picking a chinese definition, move to the next step',
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
  {
    name: 'Dicts',
    label: 'Dictionaries',
    type: DictionariesList,
  } as UserSetting,
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
  {
    name: 'ModelMappings',
    label: 'ModelManager',
    type: ModelManager,
  } as UserSetting,
  checkBox(
    'AddProgramTag',
    'Add read-chinese tag',
  ),
  checkBox(
    'AddBookTag',
    'Add source book title tag',
  ),
  checkBox(
    'AllowDuplicates',
    'Allow Duplicates',
  ),
  checkBox(
    'GenerateTermAudio',
    'Auto generate audio for keyword',
  ),
  checkBox(
    'GenerateSentenceAudio',
    'Auto generate audio for example sentence',
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
  checkBox(
    'OnlyFavorites',
    'Only show favorited books',
  ),
  checkBox(
    'HideRead',
    'Hide read books',
  ),
  checkBox(
    'ProblemFlagged',
    'Flagged Cards',
  ),
  checkBox(
    'ProblemMissingImage',
    'Missing Images',
  ),
  checkBox(
    'ProblemMissingSentence',
    'Missing Sentence',
  ),
  checkBox(
    'ProblemMissingSentenceAudio',
    'Missing Sentence Audio',
  ),
  checkBox(
    'ProblemMissingWordAudio',
    'Missing Word Audio',
  ),
  checkBox(
    'ProblemMissingPinyin',
    'Missing Pinyin',
  ),
]);

export const DisplayTable: UserConfigDisplay = {
  meta: {},
  convertValues: {},
  CardCreation: {},
  AnkiConfig: {},
  Dictionaries: {},
  SentenceGeneration: {},
  BookLibrary: {},
  CardManagement: {},
};

export function getDisplayable(obj: {[key:string]:any})
  : [any, UserSetting][] {
  return Object.entries(obj).filter(([key, _]) => {
    return key in ComponentTable;
  }).map(([key, value]) => {
    return [value, ComponentTable[key]];
  });
}

export async function generateUserSettings() {
  const settings = reactive(await GetUserSettings());

  updateSettings = async () => {
    // TODO this is a hacky way to update settings.
    // It should be set locally and not have to do these syncs
    // as they can end up having things out of sync in the front
    // end
    const newSettings = await GetUserSettings();
    Object.entries(newSettings).forEach(([key, subList]) => {
      if (key === 'convertValues') {
        return;
      }
      Object.entries(subList).forEach(([usersetting, value]) => {
        if (key === 'convertValues') {
          return;
        }
        // @ts-ignore
        const origVal = settings[key][usersetting];
        if (origVal !== value) {
          console.log('changed:', key, usersetting, value, origVal);
          // @ts-ignore
          settings[key][usersetting] = value;
        }
      });
    });
  };

  return settings;
}

type UserSettingsType = Awaited<ReturnType<typeof generateUserSettings>>;
export const UserSettingsKey = Symbol('u') as InjectionKey<UserSettingsType>;

export function getUserSettings():UserSettingsType {
  return inject(UserSettingsKey) as UserSettingsType;
}
