import { app, dialog } from 'electron';
import fs from 'fs';
import { addWord, saveLegacyWords } from '../background/knownWords';
import { importCalibreBooks } from '../background/calibre';
import { importAnkiKeywords } from '../background/ankiInterface';
import { addDictionary } from '../background/dictionaries';

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
      label: 'Add Dictionary',
      click: () => {
        const dict = dialog.showOpenDialogSync({
          properties: ['openFile'],
          filters: [
            { name: 'Yomichan Json Format', extensions: ['json'] },
          ],
        });
        addDictionary(dict[0]);
      },
    },
    {
      label: 'Import Legacy Words',
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
        saveLegacyWords(words);
      },
    },
    {
      label: 'Import Words From CSV',
      click: () => {
        // TODO handle bad selections
        const wordsFile = dialog.showOpenDialogSync({
          properties: ['openFile'],
          filters: [
            { name: 'one per line', extensions: ['csv'] },
          ],
        });
        console.log(wordsFile);
        const contents = fs.readFileSync(wordsFile[0], {
          encoding: 'utf-8',
          flags: 'r',
        });
        const words = contents.split('\n');
        words.forEach((word) => {
          addWord(word);
        });
      },
    },
    {
      label: 'Import Calibre',
      click: () => {
        const calibreDir = dialog.showOpenDialogSync({
          properties: ['openDirectory'],
        });
        console.log(calibreDir);
        importCalibreBooks(calibreDir);
      },
    },
    {
      label: 'Sync Anki',
      click: async () => {
        importAnkiKeywords();
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
