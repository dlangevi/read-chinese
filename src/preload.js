// import { ipcRenderer } from 'electron';

// window.ipcRenderer = ipcRenderer;
// window.dave = 32;

import { contextBridge, ipcRenderer } from 'electron';

// Maybe later we will want to set this up but meh
// const validChannels = ['READ_FILE', 'WRITE_FILE'];
contextBridge.exposeInMainWorld('ipc', {
  send: (channel, data) => {
    // if (validChannels.includes(channel)) {
    //  ipcRenderer.send(channel, data);
    // }
    ipcRenderer.send(channel, data);
  },
  on: (channel, func) => {
    // if (validChannels.includes(channel)) {
    //   Strip event as it includes `sender` and is a security risk
    //   ipcRenderer.on(channel, (event, ...args) => func(...args));
    // }
    console.log(channel, func);
    ipcRenderer.on(channel, (event, ...args) => {
      func(...args);
    });
  },
  loadBooks: () => ipcRenderer.invoke('loadBooks'),
  loadBook: (title) => ipcRenderer.invoke('loadBook', title),
});
