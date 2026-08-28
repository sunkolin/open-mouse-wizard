package service

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type StatusCallback func(message string)

type MouseMover struct {
	mu         sync.Mutex
	running    bool
	stopCh     chan struct{}
	minSeconds int
	maxSeconds int
}

func NewMouseMover() *MouseMover {
	return &MouseMover{
		stopCh: make(chan struct{}),
	}
}

func (m *MouseMover) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

func (m *MouseMover) Start(intervalSeconds int, onStatus StatusCallback) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("鼠标精灵已在运行中")
	}
	m.running = true
	m.minSeconds = intervalSeconds
	m.maxSeconds = intervalSeconds
	m.stopCh = make(chan struct{})
	m.mu.Unlock()

	onStatus(fmt.Sprintf("已启动，每 %d 秒移动一次鼠标", intervalSeconds))

	go m.runLoop(onStatus)

	return nil
}

func (m *MouseMover) StartRange(minSeconds, maxSeconds int, onStatus StatusCallback) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("鼠标精灵已在运行中")
	}
	m.running = true
	m.minSeconds = minSeconds
	m.maxSeconds = maxSeconds
	m.stopCh = make(chan struct{})
	m.mu.Unlock()

	onStatus(fmt.Sprintf("已启动，每 %d-%d 秒随机移动鼠标", minSeconds, maxSeconds))

	go m.runLoop(onStatus)

	return nil
}

func (m *MouseMover) runLoop(onStatus StatusCallback) {
	for {
		interval := m.randomInterval()
		timer := time.NewTimer(time.Duration(interval) * time.Second)

		select {
		case <-timer.C:
			err := moveMouseRandomly()
			timestamp := time.Now().Format("15:04:05")
			if err != nil {
				onStatus(fmt.Sprintf("[%s] 移动失败: %v", timestamp, err))
			} else {
				onStatus(fmt.Sprintf("[%s] 鼠标已移动，下次在 %d 秒后", timestamp, interval))
			}
		case <-m.stopCh:
			timer.Stop()
			return
		}
	}
}

func (m *MouseMover) randomInterval() int {
	if m.maxSeconds <= m.minSeconds {
		return m.minSeconds
	}
	return m.minSeconds + rand.Intn(m.maxSeconds-m.minSeconds+1)
}

func (m *MouseMover) Stop(onStatus StatusCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	m.running = false
	close(m.stopCh)

	onStatus("已停止")
}

func moveMouseRandomly() error {
	offsetX := rand.Intn(100) - 50
	offsetY := rand.Intn(100) - 50
	return moveMouse(offsetX, offsetY)
}

func moveMouse(offsetX, offsetY int) error {
	switch runtime.GOOS {
	case "darwin":
		return moveMouseDarwinImpl(offsetX, offsetY)
	case "windows":
		return moveMouseWindowsImpl(offsetX, offsetY)
	case "linux":
		return moveMouseLinux(offsetX, offsetY)
	default:
		return fmt.Errorf("不支持的平台: %s", runtime.GOOS)
	}
}

var moveMouseDarwinImpl = func(offsetX, offsetY int) error {
	return fmt.Errorf("macOS 实现未加载")
}

var moveMouseWindowsImpl = func(offsetX, offsetY int) error {
	return fmt.Errorf("Windows 实现未加载")
}

func moveMouseLinux(offsetX, offsetY int) error {
	getPos := exec.Command("xdotool", "getmouselocation")
	output, err := getPos.Output()
	if err != nil {
		return fmt.Errorf("xdotool 执行失败: %v", err)
	}

	var x, y int
	fmt.Sscanf(string(output), "x:%d y:%d", &x, &y)

	newX := x + offsetX
	newY := y + offsetY

	cmd := exec.Command("xdotool", "mousemove", strconv.Itoa(newX), strconv.Itoa(newY))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("移动鼠标失败: %v", err)
	}

	return nil
}

func runHelperBinary(helperBinary []byte, offsetX, offsetY int) error {
	tmpFile, err := os.CreateTemp("", "mouse-helper-*")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(helperBinary); err != nil {
		tmpFile.Close()
		return fmt.Errorf("写入辅助程序失败: %v", err)
	}
	tmpFile.Close()

	if err := os.Chmod(tmpPath, 0755); err != nil {
		return fmt.Errorf("设置执行权限失败: %v", err)
	}

	cmd := exec.Command(tmpPath, strconv.Itoa(offsetX), strconv.Itoa(offsetY))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("辅助程序执行失败: %v (%s)", err, string(output))
	}

	return nil
}
