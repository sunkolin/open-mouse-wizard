package service

import (
	"fmt"
	"math/rand"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"bling/pkg/logger"

	"github.com/go-vgo/robotgo"
)

type StatusCallback func(message string)

type MouseMover struct {
	mu         sync.Mutex
	running    bool
	stopCh     chan struct{}
	minSeconds int
	maxSeconds int
}

const (
	inputMouse          = 0
	mouseeventFMove     = 0x0001
	mouseeventFAbsolute = 0x8000
	mouseeventFLeftDown = 0x0002
	mouseeventFLeftUp   = 0x0004
)

type mouseInput struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	_           uint32
	DwExtraInfo uintptr
}

type input struct {
	Type uint32
	_    uint32
	Mi   mouseInput
}

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	smCxScreen = 0
	smCyScreen = 1
)

func NewMouseMover() *MouseMover {
	logger.Info("NewMouseMover created")
	logger.WriteDump()

	w, h := robotgo.GetScreenSize()
	logger.Info("Screen size(robotgo): %dx%d", w, h)

	smW, _, _ := procGetSystemMetrics.Call(smCxScreen)
	smH, _, _ := procGetSystemMetrics.Call(smCyScreen)
	logger.Info("Screen size(GetSystemMetrics): %dx%d", smW, smH)

	testX, testY := robotgo.GetMousePos()
	logger.Info("Initial mouse pos: (%d, %d)", testX, testY)

	probeMove := input{
		Type: inputMouse,
		Mi:   mouseInput{DwFlags: mouseeventFMove, Dx: 1, Dy: 1},
	}
	sent := robotgo.SendInput(1, unsafe.Pointer(&probeMove), int32(unsafe.Sizeof(input{})))
	logger.Info("SendInput probe (relative 1,1) returned: %d", sent)

	probeX, probeY := robotgo.GetMousePos()
	logger.Info("After probe move: (%d, %d)", probeX, probeY)

	robotgo.Move(testX, testY)
	robotgo.SendInput(1, unsafe.Pointer(&input{Type: inputMouse, Mi: mouseInput{DwFlags: mouseeventFLeftDown}}), int32(unsafe.Sizeof(input{})))
	robotgo.SendInput(1, unsafe.Pointer(&input{Type: inputMouse, Mi: mouseInput{DwFlags: mouseeventFLeftUp}}), int32(unsafe.Sizeof(input{})))

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

	logger.Info("已启动，每 %d 秒移动一次鼠标", intervalSeconds)
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

	logger.Info("已启动，每 %d-%d 秒随机移动鼠标", minSeconds, maxSeconds)
	onStatus(fmt.Sprintf("已启动，每 %d-%d 秒随机移动鼠标", minSeconds, maxSeconds))

	go m.runLoop(onStatus)

	return nil
}

func (m *MouseMover) runLoop(onStatus StatusCallback) {
	logger.Info("runLoop started")
	for {
		interval := m.randomInterval()
		timer := time.NewTimer(time.Duration(interval) * time.Second)

		select {
		case <-timer.C:
			err := moveMouseRandomly()
			timestamp := time.Now().Format("15:04:05")
			if err != nil {
				msg := fmt.Sprintf("[%s] 移动失败: %v", timestamp, err)
				logger.Error("%s", msg)
				onStatus(msg)
			} else {
				msg := fmt.Sprintf("[%s] 鼠标已移动，下次在 %d 秒后", timestamp, interval)
				logger.Info("%s", msg)
				onStatus(msg)
			}
		case <-m.stopCh:
			timer.Stop()
			logger.Info("runLoop stopped")
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

	logger.Info("Stop called")
	onStatus("已停止")
}

func moveMouseRandomly() error {
	curX, curY := robotgo.GetMousePos()
	logger.Info("moveMouseRandomly: current pos=(%d,%d)", curX, curY)

	w, h := robotgo.GetScreenSize()

	offsetX := rand.Intn(61) - 30
	offsetY := rand.Intn(61) - 30

	newX := curX + offsetX
	newY := curY + offsetY

	if w > 0 && h > 0 {
		if newX < 5 {
			newX = 5
		}
		if newX > w-5 {
			newX = w - 5
		}
		if newY < 5 {
			newY = 5
		}
		if newY > h-5 {
			newY = h - 5
		}
	}

	logger.Info("  offset=(%d,%d) target=(%d,%d)", offsetX, offsetY, newX, newY)

	moveInput := input{
		Type: inputMouse,
		Mi: mouseInput{
			Dx:      int32(offsetX),
			Dy:      int32(offsetY),
			DwFlags: mouseeventFMove,
		},
	}

	clickInputs := [3]input{
		moveInput,
		{Type: inputMouse, Mi: mouseInput{DwFlags: mouseeventFLeftDown}},
		{Type: inputMouse, Mi: mouseInput{DwFlags: mouseeventFLeftUp}},
	}
	size := int32(unsafe.Sizeof(input{}))
	sent := robotgo.SendInput(3, unsafe.Pointer(&clickInputs[0]), size)
	logger.Info("  SendInput(move+click) returned: %d", sent)

	time.Sleep(10 * time.Millisecond)

	verifyX, verifyY := robotgo.GetMousePos()
	logger.Info("  after SendInput(move+click): pos=(%d,%d)", verifyX, verifyY)

	if verifyX == curX && verifyY == curY {
		logger.Warn("  relative+click did not move, trying absolute")
		moveAbsolute(newX, newY)
		time.Sleep(15 * time.Millisecond)

		absClick := [2]input{
			{Type: inputMouse, Mi: mouseInput{DwFlags: mouseeventFLeftDown}},
			{Type: inputMouse, Mi: mouseInput{DwFlags: mouseeventFLeftUp}},
		}
		robotgo.SendInput(2, unsafe.Pointer(&absClick[0]), size)
		time.Sleep(10 * time.Millisecond)

		verifyX, verifyY = robotgo.GetMousePos()
		logger.Info("  after absolute+click: pos=(%d,%d)", verifyX, verifyY)
	}

	if verifyX == curX && verifyY == curY {
		logger.Warn("  SendInput failed, falling back to SetCursorPos+Click")
		robotgo.Move(newX, newY)
		time.Sleep(10 * time.Millisecond)
		robotgo.Click()
		time.Sleep(10 * time.Millisecond)
		verifyX, verifyY = robotgo.GetMousePos()
		logger.Info("  after SetCursorPos+Click: pos=(%d,%d)", verifyX, verifyY)
	}

	if verifyX == curX && verifyY == curY {
		logger.Error("  ALL methods failed, cursor stuck at (%d,%d)", curX, curY)
		return fmt.Errorf("鼠标未移动（可能被云桌面/系统策略阻止）")
	}

	logger.Info("  SUCCESS: (%d,%d) → (%d,%d)", curX, curY, verifyX, verifyY)
	return nil
}

func moveRelative(dx, dy int) {
	moveInput := input{
		Type: inputMouse,
		Mi: mouseInput{
			Dx:      int32(dx),
			Dy:      int32(dy),
			DwFlags: mouseeventFMove,
		},
	}
	size := int32(unsafe.Sizeof(input{}))
	robotgo.SendInput(1, unsafe.Pointer(&moveInput), size)
}

func moveAbsolute(x, y int) {
	screenW, _, _ := procGetSystemMetrics.Call(smCxScreen)
	screenH, _, _ := procGetSystemMetrics.Call(smCyScreen)

	if screenW == 0 {
		w, h := robotgo.GetScreenSize()
		screenW = uintptr(w)
		screenH = uintptr(h)
	}

	if screenW > 1 {
		absX := int32((int64(x) * 65535) / int64(screenW-1))
		absY := int32((int64(y) * 65535) / int64(screenH-1))

		moveInput := input{
			Type: inputMouse,
			Mi: mouseInput{
				Dx:      absX,
				Dy:      absY,
				DwFlags: mouseeventFMove | mouseeventFAbsolute,
			},
		}
		size := int32(unsafe.Sizeof(input{}))
		robotgo.SendInput(1, unsafe.Pointer(&moveInput), size)
	}
}
