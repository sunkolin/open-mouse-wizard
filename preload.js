const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('mouseWizard', {
  getConfig: () => ipcRenderer.invoke('get-config'),
  startMover: (interval) => ipcRenderer.invoke('start-mover', interval),
  stopMover: () => ipcRenderer.invoke('stop-mover'),
  isRunning: () => ipcRenderer.invoke('is-running'),
  onStatus: (callback) => {
    ipcRenderer.on('mover-status', (event, msg) => callback(msg));
  }
});