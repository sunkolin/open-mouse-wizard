package service

import (
	"fmt"
	"math/rand"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

type StatusCallback func(message string)

type MouseMover struct {
	mu         sync.Mutex
	running    bool
	stopCh     chan struct{}
	minSeconds int
	maxSeconds int
}

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos = user32.NewProc("GetCursorPos")
	procSetCursorPos = user32.NewProc("SetCursorPos")
)

type point struct {
	X, Y int32
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
	offsetX := rand.Intn(11) - 5
	offsetY := rand.Intn(11) - 5

	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	newX := pt.X + int32(offsetX)
	newY := pt.Y + int32(offsetY)

	procSetCursorPos.Call(uintptr(newX), uintptr(newY))
	return nil
}
