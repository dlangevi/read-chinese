import SettingsCheckbox
  from '@/components/SettingsWidgets/SettingsCheckbox.vue';
import SettingsTextbox
  from '@/components/SettingsWidgets/SettingsTextbox.vue';
import DictionariesList
  from '@/components/SettingsWidgets/DictionariesList.vue';
import SettingsSlider
  from '@/components/SettingsWidgets/SettingsSlider.vue';
import { UserSetting, UserSettingsType } from './types';

function settingsObject(
  value:string,
  label:string,
  tooltip:string,
  defaultValue:any,
  widgetType:any,
  other?:any,
):UserSetting {
  const option:UserSetting = {
    value,
    label,
    tooltip,
    defaultValue,
    type: widgetType,
    ...other,
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
    option.cached = await window.ipc.getOptionValue(value);
    option.loaded = true;
    return option.cached;
  };
  option.write = async function write(newValue:any) {
    option.cached = newValue;
    return window.ipc.setOptionValue(value, newValue);
  };
  return option;
}

function checkBox(
  value:string,
  label:string,
  tooltip:string,
  defaultValue:any,
  other?:any,
) {
  return settingsObject(
    value,
    label,
    tooltip,
    defaultValue,
    SettingsCheckbox,
    other,
  );
}

function textBox(
  value:string,
  label:string,
  tooltip:string,
  defaultValue:any,
  other?:any,
) {
  return settingsObject(
    value,
    label,
    tooltip,
    defaultValue,
    SettingsTextbox,
    other,
  );
}

function slider(
  value:string,
  label:string,
  tooltip:string,
  defaultValue:any,
  other?:any,
) {
  return settingsObject(
    value,
    label,
    tooltip,
    defaultValue,
    SettingsSlider,
    other,
  );
}

const items:UserSettingsType = (function List() {
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

  [CardCreation, Dictionaries, BookLibrary].forEach(
    async (section) => Object.values(section).forEach(
      async (option) => {
        option.readFromBackEnd();
      },
    ),
  );

  return {
    CardCreation,
    Dictionaries,
    BookLibrary,
  };
}());

export default items;
console.log(typeof items);
export type SettingsType = typeof items;
