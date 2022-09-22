import fetch from 'node-fetch';
import { addWords } from './knownWords';
import { synthesize } from './textToSpeech';
// Must have ankiconnect installed as a plugin in your anki installation
async function invoke(action, params) {
  const response = await fetch('http://localhost:8765', {
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

async function getAnkiNoteRaw(word) {
  const noteID = await invoke('findNotes', { query: `Hanzi:${word}` });
  if (noteID.result.length > 1) {
    console.log(`Too many or few notes match ${word}, ${noteID.result}`);
    return 'error';
  }
  const noteInfo = await invoke('notesInfo', {
    notes: noteID.result,
  });
  return noteInfo.result[0];
}

async function createAnkiNoteSkeleton(word) {
  // TODO for now just do this simple thing to provide a skeleton
  return {
    noteId: undefined,
    fields: {
      word,
      sentence: '',
      englishDefn: '',
      chineseDefn: '',
      pinyin: '',
      imageUrl: '',
    },
  };
}

async function getAnkiNote(word) {
  const rawNote = await getAnkiNoteRaw(word);

  // TODO make this mapping user configureable
  // to support multiple note types
  const note = {
    noteId: rawNote.noteId,
    fields: {
      word: rawNote.fields.Hanzi.value,
      sentence: rawNote.fields.ExampleSentence.value,
      englishDefn: rawNote.fields.EnglishDefinition.value,
      chineseDefn: rawNote.fields.ChineseDefinition.value,
      pinyin: rawNote.fields.Pinyin.value,
      // TODO how to load the image from a card which already exists?
      // use retrieveMediaFile to get base64 encoded image
      imageUrl: '',
    },
    rawNote,
  };
  return note;
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
async function loadFlaggedCards() {
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

async function updateCards(ankiCards) {
  addWords(ankiCards.map((ankiCard) => ({
    word: ankiCard.fields.Hanzi.value,
    interval: ankiCard.interval,
    has_flash_card: true,
  })));
}

async function updateAnkiCard(noteID, fields) {
  const res = await invoke('updateNoteFields', {
    note: {
      id: noteID,
      fields,
    },
  });
  if (res.error === null) {
    // TODO we dont always want to do this
    await removeFlag(noteID);
    return 'success';
  }
  return res;
}

export async function createAnkiCard(fields) {
  // TODO make this based on used defined fields
  const wordAudioFile = await synthesize(fields.word);
  const sentenceAudioFile = await synthesize(fields.sentence);
  console.log(fields);
  const res = await invoke('addNote', {
    note: {
      deckName: 'Reading',
      modelName: 'Reading Card',
      fields: {
        Hanzi: fields.word,
        ExampleSentence: fields.sentence,
        EnglishDefinition: fields.englishDefn,
        ChineseDefinition: fields.chineseDefn,
        Pinyin: fields.pinyin,
      },
      options: {
        allowDuplicate: true,
      },
      audio: [{
        path: wordAudioFile,
        filename: `read-chinese-hanzi-${new Date().getTime()}.wav`,
        fields: [
          'HanziAudio',
        ],
      },
      {
        path: sentenceAudioFile,
        filename: `read-chinese-sentence-${new Date().getTime()}.wav`,
        fields: [
          'SentenceAudio',
        ],
      },
      ],
      picture: [{
        url: fields.imageUrl,
        // TODO dont guess the encoding format
        filename: `read-chinese-image-${new Date().getTime()}.jpg`,
        fields: [
          'Images',
        ],
      }],
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

  updateCards(readingInfo.result);

  const skritter = await invoke('findCards', {
    query: 'deck:Skritter',
  });
  const skritterInfo = await invoke('cardsInfo', {
    cards: skritter.result,
  });
  updateCards(skritterInfo.result);
}

export const ankiInterfaceIpc = {
  createAnkiCard,
  createAnkiNoteSkeleton,
  getAnkiCard,
  getAnkiNote,
  updateAnkiCard,
  loadFlaggedCards,
  importAnkiKeywords,
};
