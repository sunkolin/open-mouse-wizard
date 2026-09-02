let robotjs;
try {
  robotjs = require('robotjs');
} catch (e) {
  robotjs = null;
}

class MouseMover {
  constructor() {
    this._running = false;
    this._timer = null;
    this._interval = 60;
    this._onStatus = null;
  }

  isRunning() {
    return this._running;
  }

  start(intervalSeconds, onStatus) {
    if (this._running) {
      throw new Error('鼠标精灵已在运行中');
    }

    this._running = true;
    this._interval = intervalSeconds;
    this._onStatus = onStatus;

    this._onStatus(`已启动，每 ${intervalSeconds} 秒点击一次Ctrl键`);

    this._scheduleNext();
  }

  _scheduleNext() {
    if (!this._running) return;

    this._timer = setTimeout(() => {
      if (!this._running) return;

      const err = this._pressCtrlKey();
      const timestamp = new Date().toLocaleTimeString('zh-CN', { hour12: false });

      if (err) {
        this._onStatus(`[${timestamp}] 点击失败: ${err.message}`);
      } else {
        this._onStatus(`[${timestamp}] Ctrl键已点击，下次在 ${this._interval} 秒后`);
      }

      this._scheduleNext();
    }, this._interval * 1000);
  }

  stop() {
    if (!this._running) return;

    this._running = false;

    if (this._timer) {
      clearTimeout(this._timer);
      this._timer = null;
    }

    if (this._onStatus) {
      this._onStatus('已停止');
    }
  }

  _pressCtrlKey() {
    try {
      if (robotjs) {
        robotjs.keyTap('control');
      } else {
        this._pressCtrlKeyFallback();
      }
      return null;
    } catch (e) {
      return e;
    }
  }

  _pressCtrlKeyFallback() {
    const { execSync } = require('child_process');
    const platform = process.platform;

    if (platform === 'darwin') {
      execSync('osascript -e \'tell application "System Events" to key code 59\'');
    } else if (platform === 'win32') {
      execSync('powershell -command "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.SendKeys]::SendWait(\'^\')"', { windowsHide: true });
    } else if (platform === 'linux') {
      execSync('xdotool key ctrl');
    } else {
      throw new Error(`不支持的平台: ${platform}`);
    }
  }
}

module.exports = MouseMover;