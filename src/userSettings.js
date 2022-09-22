import SettingsCheckbox
  from './components/SettingsWidgets/SettingsCheckbox.vue';
import SettingsTextbox
  from './components/SettingsWidgets/SettingsTextbox.vue';

function checkBox(value, label, tooltip, other) {
  const option = {
    value,
    label,
    tooltip,
    type: SettingsCheckbox,
    ...other,
    // TODO save some cached value on the renderer side
  };

  option.read = function read() {
    if (option.cached === undefined) {
      console.error(`Early read, ${option.cached}`);
    }
    return option.cached;
  };
  option.readFromBackEnd = async function readFromBackEnd() {
    option.cached = await window.ipc.getOptionValue(value);
    return option.cached;
  };
  option.write = async function write(newValue) {
    option.cached = newValue;
    return window.ipc.setOptionValue(value, newValue);
  };
  return option;
}

function textBox(value, label, tooltip, other) {
  const option = {
    value,
    label,
    tooltip,
    type: SettingsTextbox,
    ...other,
    // TODO save some cached value on the renderer side
  };

  option.read = function read() {
    if (!option.cached) {
      console.error('Early read');
    }
    return option.cached;
  };
  option.readFromBackEnd = async function readFromBackEnd() {
    option.cached = await window.ipc.getOptionValue(value);
    return option.cached;
  };
  option.write = async function write(newValue) {
    option.cached = newValue;
    return window.ipc.setOptionValue(value, newValue);
  };
  return option;
}

const items = (function List() {
  const CardCreation = {
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
    AutoAdvanceEnglish: checkBox(
      'AutoAdvanceEnglish',
      'Auto advance after definition selection',
      'After picking a definition, move to the next step',
    ),
    AutoAdvanceImage: checkBox(
      'AutoAdvanceImage',
      'Auto advance after image selection',
      'After picking a image, move to the next step',
    ),
    PopulateChinese: checkBox(
      'PopulateChinese',
      'Auto fill chinese definitions',
      'Not implemented yet',
      { disabled: true },
    ),
    GenerateTermAudio: checkBox(
      'GenerateTermAudio',
      'Auto generate audio for keyword',
      'Not implemented yet',
      { disabled: true },
    ),
    GenerateSentenceAudio: checkBox(
      'GenerateSentenceAudio',
      'Auto generate audio for example sentence',
      'Not implemented yet',
      { disabled: true },

    ),
    AutoAdvanceCard: checkBox(
      'AutoAdvanceCard',
      'Create card once all fields have been filled',
      'Create card once all fields have been filled',
    ),
  };
  const Dictionaries = {
    ShowDefinitions: checkBox(
      'ShowDefinitions',
      'Show Definitions',
      'Show the definitions for words in various tables',
    ),
    EnableChinese: checkBox(
      'EnableChinese',
      'Use Chinese definitions',
      'Allow flashcards to use chinese '
      + 'definitions instead of just english ones',
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
  };

  [CardCreation, Dictionaries].forEach(
    async (section) => Object.values(section).forEach(
      async (option) => {
        option.readFromBackEnd();
      },
    ),
  );

  return {
    CardCreation,
    Dictionaries,
  };
}());

export default items;
