const items = {
  AutoAdvanceSentence: {
    value: 'AutoAdvanceSentence',
    label: 'Auto advance after sentence selection',
    tooltip: 'After picking a sentence, move to the next step',
  },
  PopulateEnglish: {
    value: 'PopulateEnglish',
    label: 'Auto fill english definitions',
    tooltip: 'If only one definition exists, auto select it',
  },
  AutoAdvanceEnglish: {
    value: 'AutoAdvanceEnglish',
    label: 'Auto advance after definition selection',
    tooltip: 'After picking a definition, move to the next step',
  },
  PopulateChinese: {
    value: 'PopulateChinese',
    label: 'Auto fill chinese definitions',
    tooltip: 'Not implemented yet',
    disabled: true,
  },
  EnableChinese: {
    value: 'EnableChinese',
    label: 'Use Chinese definitions',
    tooltip: 'Allow flashcards to use chinese '
    + 'definitions instead of just english ones',
    disabled: true,
  },
  GenerateTermAudio: {
    value: 'GenerateTermAudio',
    label: 'Auto generate audio for keyword',
    tooltip: 'Not implemented yet',
    disabled: true,
  },
  GenerateSentenceAudio: {
    value: 'GenerateSentenceAudio',
    label: 'Auto generate audio for example sentence',
    tooltip: 'Not implemented yet',
    disabled: true,
  },
  AutoAdvanceCard: {
    value: 'AutoAdvanceCard',
    label: 'Create card once all fields have been filled',
  },
};

// TODO save some cached value on the renderer side
Object.values(items).forEach((option) => {
  option.read = async function read() {
    return window.ipc.getOptionValue(option.value);
  };
  option.write = async function write(value) {
    return window.ipc.setOptionValue(option.value, value);
  };
});

export default items;
