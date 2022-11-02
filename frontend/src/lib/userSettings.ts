import { inject } from 'vue';
import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import SettingsTextbox
  from '@/components/SettingsWidgets/SettingsTextbox.vue';
import DictionariesList
  from '@/components/SettingsWidgets/DictionariesList.vue';
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
import { UserSettingsKey, UserSetting, UserSettingsType } from './types';

function settingsObject(
  value:string,
  label:string,
  tooltip:string,
  defaultValue:any,
  widgetType:any,
  getter:any,
  setter:any,
):UserSetting {
  const option:UserSetting = {
    value,
    label,
    tooltip,
    defaultValue,
    type: widgetType,
    loaded: false,
    // TODO save some cached value on the renderer side
  };

  option.read = function read() {
    if (!option.loaded) {
      console.error(`Early read, ${option.value}`);
      return option.defaultValue;
    }
    if (option.cached === undefined) {
      return option.defaultValue;
    }
    return option.cached;
  };
  option.readFromBackEnd = async function readFromBackEnd() {
    option.cached = await getter(value, defaultValue);
    option.loaded = true;
    return option.cached;
  };
  option.write = async function write(newValue:any) {
    option.cached = newValue;
    return setter(value, newValue);
  };
  return option;
}

function checkBox(
  value:string,
  label:string,
  tooltip:string,
  defaultValue:any,
) {
  return settingsObject(
    value,
    label,
    tooltip,
    defaultValue,
    SettingsCheckbox,
    GetUserSettingBool,
    SetUserSettingBool,
  );
}

function textBox(
  value:string,
  label:string,
  tooltip:string,
  defaultValue:any,
) {
  return settingsObject(
    value,
    label,
    tooltip,
    defaultValue,
    SettingsTextbox,
    GetUserSetting,
    SetUserSetting,
  );
}

function slider(
  value:string,
  label:string,
  tooltip:string,
  defaultValue:any,
) {
  return settingsObject(
    value,
    label,
    tooltip,
    defaultValue,
    SettingsSlider,
    GetUserSettingInt,
    SetUserSettingInt,

  );
}

export async function generateUserSettings() :Promise<UserSettingsType> {
  const CardCreation = {
    AutoAdvanceSentence: checkBox(
      'AutoAdvanceSentence',
      'Auto advance after sentence selection',
      'After picking a sentence, move to the next step',
      true,
    ),
    PopulateEnglish: checkBox(
      'PopulateEnglish',
      'Auto fill english definitions',
      'If only one definition exists, auto select it',
      false,
    ),
    PopulateChinese: checkBox(
      'PopulateChinese',
      'Auto fill chinese definitions',
      'If only one definition exists, auto select it',
      false,
    ),
    AutoAdvanceEnglish: checkBox(
      'AutoAdvanceEnglish',
      'Auto advance after definition selection',
      'After picking a definition, move to the next step',
      false,
    ),
    AutoAdvanceImage: checkBox(
      'AutoAdvanceImage',
      'Auto advance after image selection',
      'After picking a image, move to the next step',
      false,
    ),
    GenerateTermAudio: checkBox(
      'GenerateTermAudio',
      'Auto generate audio for keyword',
      'Not implemented yet',
      false,
    ),
    GenerateSentenceAudio: checkBox(
      'GenerateSentenceAudio',
      'Auto generate audio for example sentence',
      'Not implemented yet',
      false,

    ),
    AutoAdvanceCard: checkBox(
      'AutoAdvanceCard',
      'Create card once all fields have been filled',
      'Create card once all fields have been filled',
      true,
    ),
  };
  const Dictionaries = {
    Dictionaries: {
      value: 'Dictionaries',
      label: 'Dictionaries',
      defaultValue: [],
      type: DictionariesList,
      readFromBackEnd: () => {},
    },
    ShowDefinitions: checkBox(
      'ShowDefinitions',
      'Show Definitions',
      'Show the definitions for words in various tables',
      true,
    ),
    EnableChinese: checkBox(
      'EnableChinese',
      'Use Chinese definitions',
      'Allow flashcards to use chinese '
      + 'definitions instead of just english ones',
      true,
    ),
    AzureApiKey: textBox(
      'AzureApiKey',
      'Azure Audio Api Key',
      'Setup an free azure tts account and put your key here',
      '',
    ),
    AzureImageApiKey: textBox(
      'AzureImageApiKey',
      'Azure Image Api Key',
      'Setup an free azure bing search and put your key here',
      '',
    ),
    KnownInterval: slider(
      'KnownInterval',
      'Time before a word is considered "known"',
      'How long of an interval in anki before a word is '
      + ' included in generated sentences',
      10,

    ),
  };
  const BookLibrary = {
    OnlyFavorites: checkBox(
      'OnlyFavorites',
      'Only show favorited books',
      '',
      false,
    ),
  };

  await Promise.all(
    [CardCreation, Dictionaries, BookLibrary].map(
      async (section) => Promise.all(
        Object.values(section).map(
          async (option) => option.readFromBackEnd(),
        ),
      ),
    ),
  );

  return {
    CardCreation,
    Dictionaries,
    BookLibrary,
  };
}

export function getUserSettings():UserSettingsType {
  return inject(UserSettingsKey) as UserSettingsType;
}
