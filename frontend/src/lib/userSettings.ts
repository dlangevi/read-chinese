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
import { backend } from '@wailsjs/models';
import { EventsOn } from '../../wailsjs/runtime/runtime';

import {
  SetUserSetting,
  SetUserSettingBool,
  SetUserSettingInt,
  GetUserSettings,
} from '@wailsjs/backend/UserSettings';

type SettingsSetter<Type> = (arg1: string, arg2: Type) => Promise<void>

type SettingsWidgets =
  typeof SettingsCheckbox |
  typeof SettingsTextbox |
  typeof SettingsSelector |
  typeof DictionariesList |
  typeof ModelManager |
  typeof SettingsSlider

export type UserSetting<T> = {
  name:string;
  label:string;
  tooltip?:string;
  type?: SettingsWidgets;
  write: (arg: T) => Promise<void>;
  dataSource?: () => Promise<string[]>;
};

function settingsObject<Type>(
  name:string,
  label:string,
  tooltip:string|undefined,
  widgetType:SettingsWidgets,
  setter:SettingsSetter<Type>,
):UserSetting<Type> {
  const option:UserSetting<Type> = {
    name,
    label,
    tooltip,
    type: widgetType,
    write: async function write(newValue: Type) {
      await setter(name, newValue);
    },
  };
  return option;
}

function selector(
  name:string,
  label:string,
  dataSource: () => Promise<string[]>,
  tooltip?:string,
):UserSetting<string> {
  const option:UserSetting<string> = {
    name,
    label,
    tooltip,
    type: SettingsSelector,
    write: async function write(newValue:string) {
      await SetUserSetting(name, newValue);
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

type UserConfigDisplay = {
  [K in keyof backend.UserSettings]: {
    [P in keyof backend.UserSettings[K]]?:
    UserSetting<backend.UserSettings[K][P]>
  }
}

export const ComponentTable : UserConfigDisplay = {
  meta: {
    EnableExperimental: checkBox(
      'EnableExperimental',
      'Enable Experimental Features',
      'Only do this if you are me',
    ),
    Theme: selector(
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
  },
  CardCreation: {
    AutoAdvanceSentence: checkBox(
      'AutoAdvanceSentence',
      'Auto advance after sentence selection',
      'After picking a sentence, move to the next step',
    ),
    PopulateEnglish: checkBox(
      'PopulateEnglish',
      'Auto fill english definitions',
      'If only one definition exists, auto select it',
    ),
    PopulateChinese: checkBox(
      'PopulateChinese',
      'Auto fill chinese definitions',
      'If only one definition exists, auto select it',
    ),
    AutoAdvanceEnglish: checkBox(
      'AutoAdvanceEnglish',
      'Auto advance after english definition selection',
      'After picking an english definition, move to the next step',
    ),
    AutoAdvanceChinese: checkBox(
      'AutoAdvanceChinese',
      'Auto advance after chinese definition selection',
      'After picking a chinese definition, move to the next step',
    ),
    AutoAdvanceImage: checkBox(
      'AutoAdvanceImage',
      'Auto advance after image selection',
      'After picking a image, move to the next step',
    ),
    AutoAdvanceCard: checkBox(
      'AutoAdvanceCard',
      'Create card once all fields have been filled',
      'Create card once all fields have been filled',
    ),
  },
  AnkiConfig: {
    ActiveDeck: selector(
      'ActiveDeck',
      'Active Anki Deck',
      LoadDecks,
    ),
    ActiveModel: selector(
      'ActiveModel',
      'Active Anki Model',
      LoadModels,
    ),
    ModelMappings: {
      name: 'ModelMappings',
      label: 'ModelManager',
      type: ModelManager,
    } as UserSetting<{ [key: string]: backend.FieldsMapping }>,
    AllowDuplicates: checkBox(
      'AllowDuplicates',
      'Allow Duplicates',
    ),
    GenerateTermAudio: checkBox(
      'GenerateTermAudio',
      'Auto generate audio for keyword',
    ),
    GenerateSentenceAudio: checkBox(
      'GenerateSentenceAudio',
      'Auto generate audio for example sentence',
    ),
    AzureApiKey: textBox(
      'AzureApiKey',
      'Azure Audio Api Key',
      'Setup an free azure tts account and put your key here',
    ),
    AzureImageApiKey: textBox(
      'AzureImageApiKey',
      'Azure Image Api Key',
      'Setup an free azure bing search and put your key here',
    ),
    AddProgramTag: checkBox(
      'AddProgramTag',
      'Add read-chinese tag',
    ),
    AddBookTag: checkBox(
      'AddBookTag',
      'Add source book title tag',
    ),
  },
  Dictionaries: {
    Dicts: {
      name: 'Dicts',
      label: 'Dictionaries',
      type: DictionariesList,
    } as UserSetting<{ [key: string]: backend.Dict }>,
    ShowDefinitions: checkBox(
      'ShowDefinitions',
      'Show Definitions',
      'Show the definitions for words in various tables',
    ),
    EnableChinese: checkBox(
      'EnableChinese',
      'Use Chinese definitions',
      'Allow flashcards to use chinese ' +
      'definitions instead of just english ones',
    ),
  },
  SentenceGeneration: {
    KnownInterval: slider(
      'KnownInterval',
      'Time before a word is considered "known"',
      'How long of an interval in anki before a word is ' +
      ' included in generated sentences',
    ),
    IdealSentenceLength: slider(
      'IdealSentenceLength',
      'Ideal Sentence Length"',
      'What the ideal sentence length you want to be selected from books',
    ),
  },
  BookLibrary: {
    OnlyFavorites: checkBox(
      'OnlyFavorites',
      'Only show favorited books',
    ),
    HideRead: checkBox(
      'HideRead',
      'Hide read books',
    ),
  },
  // TODO filter this from the type
  convertValues: {},

};

export async function generateUserSettings() {
  const settings = reactive(await GetUserSettings());

  // Still not 100% happy with this
  EventsOn('UpdatedConfig', (newSettings: backend.UserSettings) => {
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
  });

  return settings;
}

type UserSettingsType = Awaited<ReturnType<typeof generateUserSettings>>;
export const UserSettingsKey = Symbol('u') as InjectionKey<UserSettingsType>;

export function getUserSettings():UserSettingsType {
  return inject(UserSettingsKey) as UserSettingsType;
}
