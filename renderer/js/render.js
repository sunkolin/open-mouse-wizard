const intervalInput = document.getElementById('interval');
const toggleBtn = document.getElementById('toggleBtn');
const statusBar = document.getElementById('statusBar');

let isRunning = false;

async function init() {
  try {
    const config = await window.mouseWizard.getConfig();
    if (config && config['default-interval']) {
      intervalInput.value = config['default-interval'];
    }
  } catch (e) {
    console.error('Failed to load config:', e);
  }

  window.mouseWizard.onStatus((msg) => {
    statusBar.textContent = msg;
  });
}

toggleBtn.addEventListener('click', async () => {
  if (isRunning) {
    await stopMover();
  } else {
    await startMover();
  }
});

async function startMover() {
  const interval = parseInt(intervalInput.value, 10);
  if (isNaN(interval) || interval < 1) {
    statusBar.textContent = '请输入有效的正整数';
    return;
  }

  try {
    const result = await window.mouseWizard.startMover(interval);
    if (result.success) {
      isRunning = true;
      toggleBtn.textContent = '停止';
      toggleBtn.className = 'btn btn-stop';
    } else {
      statusBar.textContent = `启动失败: ${result.error}`;
    }
  } catch (e) {
    statusBar.textContent = `启动失败: ${e.message}`;
  }
}

async function stopMover() {
  try {
    await window.mouseWizard.stopMover();
    isRunning = false;
    toggleBtn.textContent = '启动';
    toggleBtn.className = 'btn btn-start';
  } catch (e) {
    statusBar.textContent = `停止失败: ${e.message}`;
  }
}

init();