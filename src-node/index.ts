import readline from 'readline';
import express from 'express';
import bodyParser from 'body-parser';
import { syncWords, checkWords } from './knownWords';
import { loadDictionaries } from './dictionaries';
import { preloadWords } from './segmentation';
import {
  updateTimesRan,
  getTimesRan,
  initializeDatabases,
} from './database';
import { registerIpc } from './ipcLoader';

const userConfigDir = process.argv[2];
//
async function main() {
  console.time('bootup');
  console.log('init database ...');
  await initializeDatabases(userConfigDir);
  updateTimesRan();
  console.log(`Ran ${getTimesRan()}`);
  console.log('syncWords ...');
  await syncWords();
  console.log('load dictionaries ...');
  loadDictionaries();
  console.log('preload words ...');
  await preloadWords(userConfigDir);
  console.log('done ...');
  console.timeEnd('bootup');
  checkWords();
}
main();

function onMessage(cb: (ling:string) => void) {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false,
  });

  rl.on('line', (line) => {
    cb(line);
  });
}
function write(message:string) {
  console.log(message);
}

onMessage((line:string) => {
  write(`read ${line}`);
});

const app = express();
app.use(bodyParser.json());

registerIpc(app);

const port = 3451;
app.listen(port, () => {
  console.log(`Server started on port ${port}`);
});

// TODO
// electron apis I need to make wrappers for
// dialog
// resourcesPath
// userData
