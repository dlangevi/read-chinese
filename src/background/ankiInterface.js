import fetch from 'node-fetch';
import { addWord } from './knownWords';
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

export async function addSentenceToCard(word, sentence) {
  const notes = await invoke('findNotes', { query: `Hanzi:${word}` });
  if (notes.result.length !== 1) {
    console.log(`Too many or few notes match ${word}, ${notes.result}`);
    return;
  }
  const noteID = notes.result[0];
  console.log(`Adding sentence to card for ${word} (noteID ${noteID})`);
  const result = await invoke('updateNoteFields', {
    note: {
      id: noteID,
      fields: {
        ExampleSentence: sentence,
      },
    },
  });
  if (result.error) {
    console.log(result);
  }
}

async function getAnkiNote(word) {
  const noteID = await invoke('findNotes', { query: `Hanzi:${word}` });
  if (noteID.result.length !== 1) {
    console.log(`Too many or few notes match ${word}, ${noteID.result}`);
    return 'error';
  }
  const noteInfo = await invoke('notesInfo', {
    notes: noteID.result,
  });
  return noteInfo.result[0];
}

async function getAnkiCard(word) {
  const cardID = await invoke('findCards', { query: `Hanzi:${word}` });
  if (cardID.result.length !== 1) {
    console.log(`Too many or few notes match ${word}, ${cardID.result}`);
    return 'error';
  }
  const cardInfo = await invoke('cardsInfo', {
    cards: cardID.result,
  });
  return cardInfo.result[0];
}

export async function getLackingCards(deck) {
  const skritter = await invoke('findCards', {
    query: `deck:${deck} ExampleSentence:`,
  });
  const skritterInfo = await invoke('cardsInfo', {
    cards: skritter.result,
  });
  const allWords = skritterInfo.result
    .map((card) => fixWord(card.fields.Hanzi.value))
    .filter((word) => isChinese(word));
  return allWords;
}

// TODO filter by deck
async function getFlaggedCards() {
  const flaggedIDs = await invoke('findCards', {
    query: 'flag:1',
  });
  const flaggedCards = await invoke('cardsInfo', {
    cards: flaggedIDs.result,
  });
  return flaggedCards.result
    .map((card) => ({
      word: fixWord(card.fields.Hanzi.value),
      sentence: card.fields.ExampleSentence.value,
    }));
}

async function updateCard(ankiCard) {
  addWord(
    ankiCard.fields.Hanzi.value,
    ankiCard.interval,
    true,
  );
}

async function updateAnkiCard(noteID, fields) {
  const res = await invoke('updateNoteFields', {
    note: {
      id: noteID,
      fields,
    },
  });
  if (res.error === null) {
    return 'success';
  }
  return res;
}

async function removeFlag(noteID) {
  const noteInfo = await invoke('notesInfo', {
    notes: [noteID],
  });
  const note = noteInfo.result[0];
  // TODO handle this
  if (note.cards.length !== 1) {
    console.error('Note has multiple cards');
  }
  await invoke(
    'setSpecificValueOfCard',
    {
      card: note.cards[0],
      keys: ['flags'],
      newValues: [0],
    },
  );
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
export function initAnkiIpc(ipcMain) {
  ipcMain.handle('getAnkiCard', async (event, word) => getAnkiCard(word));
  ipcMain.handle('getAnkiNote', async (event, word) => getAnkiNote(word));
  ipcMain.handle('updateAnkiCard', async (event, noteID, fields) => {
    await removeFlag(noteID);
    return updateAnkiCard(noteID, fields);
  });
  ipcMain.handle('flaggedCards', async () => {
    const flagged = await getFlaggedCards('Reading');
    return flagged;
  });
}
