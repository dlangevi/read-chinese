import { app, dialog } from 'electron';
import { importCalibreBooks } from '../calibre';
import { addDictionary } from '../dictionaries';

export default {
  label: 'App',
  submenu: [
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
      label: 'Quitter',
      accelerator: 'CmdOrCtrl+Q',
      click: () => {
        app.quit();
      },
    },
  ],
};
