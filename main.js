const { app, BrowserWindow, Tray, Menu, ipcMain, nativeImage } = require('electron');
const path = require('path');
const MouseMover = require('./service/mouse-mover');

let mainWindow = null;
let tray = null;
let mover = new MouseMover();

function loadConfig() {
  const fs = require('fs');

  const defaultConfig = {
    title: '鼠标精灵',
    width: 200,
    height: 200,
    'default-interval': 60
  };

  const candidates = [
    path.join(app.getAppPath(), 'config', 'config.json'),
    path.join(app.getAppPath(), 'config.json'),
    path.join(process.resourcesPath || '', 'config', 'config.json'),
    path.join(path.dirname(app.getPath('exe')), 'config', 'config.json')
  ];

  for (const configPath of candidates) {
    try {
      if (fs.existsSync(configPath)) {
        const data = fs.readFileSync(configPath, 'utf-8');
        const config = JSON.parse(data);
        return { ...defaultConfig, ...config };
      }
    } catch (e) {
      continue;
    }
  }

  return defaultConfig;
}

function createWindow() {
  const config = loadConfig();

  mainWindow = new BrowserWindow({
    width: config.width,
    height: config.height,
    resizable: false,
    title: config.title,
    icon: getAppIcon(),
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false
    }
  });

  mainWindow.loadFile(path.join(__dirname, 'renderer', 'index.html'));

  mainWindow.on('closed', () => {
    mainWindow = null;
  });

  mainWindow.on('close', (event) => {
    if (tray && !app.isQuitting) {
      event.preventDefault();
      mainWindow.hide();
    }
  });
}

function getAppIcon() {
  const iconPath = path.join(app.getAppPath(), 'assets', 'icon.png');
  try {
    return nativeImage.createFromPath(iconPath);
  } catch (e) {
    return undefined;
  }
}

function createTray() {
  const icon = getAppIcon();
  if (!icon || icon.isEmpty()) return;

  tray = new Tray(icon);
  const contextMenu = Menu.buildFromTemplate([
    { label: '显示主界面', click: () => { if (mainWindow) mainWindow.show(); } },
    { type: 'separator' },
    { label: '退出', click: () => { app.isQuitting = true; app.quit(); } }
  ]);
  tray.setToolTip('鼠标精灵');
  tray.setContextMenu(contextMenu);
  tray.on('double-click', () => { if (mainWindow) mainWindow.show(); });
}

function setupIPC() {
  ipcMain.handle('get-config', () => {
    return loadConfig();
  });

  ipcMain.handle('start-mover', (event, interval) => {
    try {
      mover.start(interval, (msg) => {
        if (mainWindow && !mainWindow.isDestroyed()) {
          mainWindow.webContents.send('mover-status', msg);
        }
      });
      return { success: true };
    } catch (e) {
      return { success: false, error: e.message };
    }
  });

  ipcMain.handle('stop-mover', () => {
    mover.stop();
    return { success: true };
  });

  ipcMain.handle('is-running', () => {
    return mover.isRunning();
  });
}

app.whenReady().then(() => {
  createWindow();
  createTray();
  setupIPC();
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', () => {
  mover.stop();
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  } else if (mainWindow) {
    mainWindow.show();
  }
});