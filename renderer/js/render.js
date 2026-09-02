const intervalInput = document.getElementById('interval');
const toggleBtn = document.getElementById('toggleBtn');

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
    return;
  }

  try {
    const result = await window.mouseWizard.startMover(interval);
    if (result.success) {
      isRunning = true;
      toggleBtn.textContent = '停止';
      toggleBtn.className = 'btn btn-stop';
    }
  } catch (e) {
    console.error('Start failed:', e.message);
  }
}

async function stopMover() {
  try {
    await window.mouseWizard.stopMover();
    isRunning = false;
    toggleBtn.textContent = '启动';
    toggleBtn.className = 'btn btn-start';
  } catch (e) {
    console.error('Stop failed:', e.message);
  }
}

init();