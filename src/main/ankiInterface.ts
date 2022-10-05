import fetch from 'node-fetch';
import { addWords, addWord } from './knownWords';
import { synthesize } from './textToSpeech';
import { getOptionValue } from './database';
// Must have ankiconnect installed as a plugin in your anki installation
async function invoke(action:string, params:any):Promise<any> {
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

type Fields = {
  word: string,
  sentence: string,
  englishDefn: string,
  chineseDefn: string,
  pinyin: string,
  imageUrls: [string],
}

// eslint-disable-next-line no-unused-vars, @typescript-eslint/no-unused-vars
function isChinese(word:string) {
  // unicode ranges for chinese characters
  const isOnlyChinese = /^[\u3040-\u30ff\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff\uff66-\uff9f]*$/
    .test(word);
  if (!isOnlyChinese) {
    console.log(`${word} is sus, skipping it`);
  }
  return isOnlyChinese;
}

function fixWord(origWord:string) {
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

async function getAnkiNoteRaw(word:string) {
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

async function createAnkiNoteSkeleton(word:string) {
  // TODO for now just do this simple thing to provide a skeleton
  return {
    noteId: undefined,
    fields: {
      word,
      sentence: '',
      englishDefn: '',
      chineseDefn: '',
      pinyin: '',
      imageUrls: [],
    },
  };
}

async function getAnkiNote(word:string) {
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
      imageUrls: [],
    },
    rawNote,
  };
  return note;
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
    .map((card:any) => ({
      word: fixWord(card.fields.Hanzi.value),
      sentence: card.fields.ExampleSentence.value,
    }));
}

async function updateCards(ankiCards:any) {
  addWords(ankiCards.map((ankiCard:any) => ({
    word: ankiCard.fields.Hanzi.value,
    interval: ankiCard.interval,
    has_flash_card: true,
  })));
}

// TODO use Fields object
async function updateAnkiCard(noteID:any, fields:any) {
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

async function createAnkiCard(fields:Fields, tags = []) {
  // TODO make this based on used defined fields
  console.log('tags', tags);
  const audioArray = [];
  if (getOptionValue('GenerateTermAudio', false)) {
    const wordAudioFile = await synthesize(fields.word);
    audioArray.push({
      path: wordAudioFile,
      filename: `read-chinese-hanzi-${new Date().getTime()}.wav`,
      fields: [
        'HanziAudio',
      ],
    });
  }
  if (getOptionValue('GenerateSentenceAudio', false)) {
    const sentenceAudioFile = await synthesize(fields.sentence);
    audioArray.push({
      path: sentenceAudioFile,
      filename: `read-chinese-sentence-${new Date().getTime()}.wav`,
      fields: [
        'SentenceAudio',
      ],
    });
  }
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
      tags: [
        'read-chinese',
        ...tags,
      ],
      options: {
        allowDuplicate: true,
      },
      audio: audioArray,
      picture:
        fields.imageUrls.map((imageUrl:string) => ({
          url: imageUrl,
          // TODO dont guess the encoding format
          filename: `read-chinese-image-${new Date().getTime()}.jpg`,
          fields: [
            'Images',
          ],
        })),
    },
  });
  if (res.error === null) {
    addWord(fields.word, 0, true);
    return 'success';
  }
  return res;
}

async function removeFlag(noteID:string) {
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
async function importAnkiKeywords() {
  const reading = await invoke('findCards', {
    query: 'deck:Reading',
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
  getAnkiNote,
  updateAnkiCard,
  loadFlaggedCards,
  importAnkiKeywords,
};
