const items = {
  PopulateEnglish: {
    value: 'PopulateEnglish',
    label: 'Auto fill english definitions',
    tooltip: 'If only one definition exists, auto select it',
  },
  PopulateChinese: {
    value: 'PopulateChinese',
    label: 'Auto fill chinese definitions',
    tooltip: 'Not implemented yet',
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
};

Object.values(items).forEach((option) => {
  option.read = async function read() {
    return window.ipc.getOptionValue(option.value);
  };
  option.write = async function write(value) {
    return window.ipc.setOptionValue(option.value, value);
  };
});

export default items;
