import { app, dialog } from 'electron';
import { createParser } from 'node-csv';
import fs from 'fs';
import { saveWords } from '../helpers/database';
import { importCalibreBooks } from '../helpers/calibre';
import { generateSentences } from '../helpers/generateSentences';
import { getSkritterWords } from '../helpers/ankiInterface';

const csv = createParser();

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
      label: 'Import Calibre',
      click: () => {
        importCalibreBooks();
      },
    },
    {
      label: 'Test Sentence Gen',
      click: () => {
        const wordsFile = dialog.showOpenDialogSync({
          properties: ['openFile'],
          filters: [
            { name: 'csv', extensions: ['csv'] },
          ],
        });
        console.log(wordsFile);
        const contents = fs.readFileSync(wordsFile[0], {
          encoding: 'utf-8',
          flags: 'r',
        });
        // TODO load csv words and generate known sentences
        const words = csv.parse(contents);
        const fixedWords = words
          .filter((row) => (row.length === 4))
          .map((row) => row[0]);

        generateSentences(fixedWords);
      },
    },
    {
      label: 'Test Anki Gen',
      click: async () => {
        const ankiWords = await getSkritterWords();
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
