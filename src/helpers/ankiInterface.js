import fetch from 'node-fetch';
import { updateWord } from './database';
// Must have ankiconnect installed as a plugin in your anki installation
async function invoke(action, params) {
  const response = await fetch('http://localhost:8765', {
  // const response = await fetch('http://120.0.0.1:8765', {
    method: 'Post',
    body: JSON.stringify({
      action,
      version: 6,
      params: {
        ...params,
      },
    }),
  });
  return response.json();
}

function isChinese(word) {
  // unicode ranges for chinese characters
  const isOnlyChinese = /^[\u3040-\u30ff\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff\uff66-\uff9f]*$/
    .test(word);
  if (!isOnlyChinese) {
    console.log(`${word} is sus, skipping it`);
  }
  return isOnlyChinese;
}

function fixWord(origWord) {
  let word = origWord;
  word = word.replace(/<br>/gi, '');
  word = word.replace(/<div>/gi, '');
  word = word.replace(/<\/div>/gi, '');
  word = word.replace(/,/gi, '');
  word = word.replace(/&nbsp/, '');

  if (word !== origWord) {
    console.log(`warning ${word} is ${origWord}`);
  }
  return word;
}

export async function createCard(/* card */) {
  // const result = await invoke('modelNames', {
  const result = await invoke('addNote', {
    note: {
      deckName: 'Testing',
      modelName: 'Reading Card',
      fields: {
        Simplified: 'test',
        Meaning: 'test',
        EnglishMeaning: 'test',
        SentenceSimplified: 'test',
      },
      options: {
        allowDuplicate: true,
      },
      audio: [{
        url: 'https://assets.languagepod101.com/dictionary/japanese/audiomp3.php?kanji=猫&kana=ねこ',
        filename: 'yomichan_ねこ_猫.mp3',
        skipHash: '7e2c2f954ef6051373ba916f000168dc',
        fields: [
          'Audio',
        ],
      }],
      picture: [{
        url: 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c7/A_black_cat_named_Tilly.jpg/220px-A_black_cat_named_Tilly.jpg',
        filename: 'black_cat.jpg',
        skipHash: '8d6e4646dfae812bf39651b59d7429ce',
        fields: [
          'SentenceImage',
        ],
      }],
    },
  });
  return result;
}

export async function getSkritterWords() {
  const skritter = await invoke('findCards', { query: 'deck:Skritter' });
  const skritterInfo = await invoke('cardsInfo', {
    cards: skritter.result,
  });
  const allWords = skritterInfo.result
    .map((card) => fixWord(card.fields.Simplified.value))
    .filter((word) => isChinese(word));
  return allWords;
}

async function updateCard(ankiCard) {
  updateWord(ankiCard);
}

// @todo: make this configurable from the app to pick certian decks and fields
export async function importAnkiKeywords() {
  const reading = await invoke('findCards', {
    query: 'deck:Reading -"note:Audio Card"',
  });
  const readingInfo = await invoke('cardsInfo', {
    cards: reading.result,
  });

  readingInfo.result.forEach((card) => updateCard(card));

  const skritter = await invoke('findCards', {
    query: 'deck:Skritter',
  });
  const skritterInfo = await invoke('cardsInfo', {
    cards: skritter.result,
  });
  skritterInfo.result.forEach((card) => updateCard(card));
}
