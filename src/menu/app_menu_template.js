import { app, dialog } from 'electron';
import fs from 'fs';
import { saveWords, loadWords } from '../helpers/database';
import { importCalibreBooks } from '../helpers/calibre';
import { generateSentences } from '../helpers/generateSentences';
import { getLackingCards, importAnkiKeywords } from '../helpers/ankiInterface';

export default {
  label: 'App',
  submenu: [
    {
      label: 'Add Book',
      click: () => {
        console.log(dialog.showOpenDialogSync({
          properties: ['openFile'],
          filters: [
            { name: 'Plain Text', extensions: ['txt'] },
          ],
        }));
      },
    },
    {
      label: 'Import Words',
      click: () => {
        // TODO handle bad selections
        const wordsFile = dialog.showOpenDialogSync({
          properties: ['openFile'],
          filters: [
            { name: 'Json format', extensions: ['json'] },
          ],
        });
        console.log(wordsFile);
        const contents = fs.readFileSync(wordsFile[0], {
          encoding: 'utf-8',
          flags: 'r',
        });
        const words = JSON.parse(contents);
        saveWords(words);
      },
    },
    {
      label: 'Log Words',
      click: async () => {
        console.log('before');
        const words = await loadWords();
        console.log(`after ${words.length}`);
      },
    },
    {
      label: 'Import Calibre',
      click: () => {
        importCalibreBooks();
      },
    },
    {
      label: 'Sync Anki',
      click: async () => {
        importAnkiKeywords();
      },
    },
    {
      label: 'Auto Generate Missing Sentences',
      click: async () => {
        const ankiWords = await getLackingCards('Reading');
        console.log(ankiWords);
        generateSentences(ankiWords);
      },
    },
    {
      label: 'Quitter',
      accelerator: 'CmdOrCtrl+Q',
      click: () => {
        app.quit();
      },
    },
  ],
};
