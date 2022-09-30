// The built directory structure
//
// ├─┬ dist
// │ ├─┬ electron
// │ │ ├─┬ main
// │ │ │ └── index.js
// │ │ └─┬ preload
// │ │   └── index.js
// │ ├── index.html
// │ ├── ...other-static-files-from-public
// │
import {
  app, BrowserWindow, shell, protocol,
} from 'electron';
import { release } from 'os';
import { join } from 'path';

import { syncWords, checkWords } from './background/knownWords';
import { loadDictionaries } from './background/dictionaries';
import { preloadWords } from './background/segmentation';
import { initIpcMain } from '../ipcLoader';
import {
  updateTimesRan,
  getTimesRan,
  initializeDatabase,
} from './background/database';

console.time('bootup');
process.env.DIST = join(__dirname, '../..');
process.env.PUBLIC = app.isPackaged
  ? process.env.DIST : join(process.env.DIST, '../public');
console.log(process.resourcesPath);

// Disable GPU Acceleration for Windows 7
if (release().startsWith('6.1')) app.disableHardwareAcceleration();

// Set application name for Windows 10+ notifications
if (process.platform === 'win32') app.setAppUserModelId(app.getName());

if (!app.requestSingleInstanceLock()) {
  app.quit();
  process.exit(0);
}

// Remove electron security warnings
// This warning only shows in development mode
// Read more on https://www.electronjs.org/docs/latest/tutorial/security
// process.env['ELECTRON_DISABLE_SECURITY_WARNINGS'] = 'true'

let win: BrowserWindow | null = null;
// Here, you can also use other preload
const preload = join(__dirname, '../preload/index.js');
const devUrl = process.env.VITE_DEV_SERVER_URL as string;
const indexHtml = join(process.env.DIST, 'index.html');

async function createWindow() {
  updateTimesRan();
  console.log(`Ran ${getTimesRan()}`);
  console.log(`Data stored at ${app.getPath('userData')}`);
  win = new BrowserWindow({
    title: 'read-chinese',
    icon: join(process.env.PUBLIC, 'favicon.ico'),
    webPreferences: {
      preload,
      nodeIntegration: false,
      sandbox: false,
      contextIsolation: true,
    },
  });

  if (app.isPackaged) {
    win.loadFile(indexHtml);
  } else {
    win.loadURL(devUrl);
    // Open devTool if the app is not packaged
    win.webContents.openDevTools();
  }

  // Test actively push message to the Electron-Renderer
  win.webContents.on('did-finish-load', () => {
    win?.webContents.send('main-process-message', new Date().toLocaleString());
  });

  // Make all links open with the browser, not with the application
  win.webContents.setWindowOpenHandler(({ url }) => {
    if (url.startsWith('https:')) shell.openExternal(url);
    return { action: 'deny' };
  });
}

app.whenReady().then(async () => {
  protocol.interceptFileProtocol(
    'atom',
    (request, callback) => {
      const pathname = request.url.replace('atom:///', '');
      callback(pathname);
    },
  );
  console.log('init database ...');
  await initializeDatabase();
  console.log('syncWords ...');
  await syncWords();
  // setApplicationMenu();
  console.log('init IPC...');
  initIpcMain();
  console.log('create Window...');
  await createWindow();
  console.log('load dictionaries ...');
  await loadDictionaries();
  console.log('preload words ...');
  await preloadWords();
  console.log('done ...');
  console.timeEnd('bootup');
  checkWords();
});

app.on('window-all-closed', () => {
  win = null;
  if (process.platform !== 'darwin') app.quit();
});

app.on('second-instance', () => {
  if (win) {
    // Focus on the main window if the user tried to open another
    if (win.isMinimized()) win.restore();
    win.focus();
  }
});

app.on('activate', async () => {
  const allWindows = BrowserWindow.getAllWindows();
  if (allWindows.length) {
    allWindows[0].focus();
  } else {
    createWindow();
  }
});
