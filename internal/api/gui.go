package api

import (
	"bling/internal/config"
	"bling/internal/model"
	"bling/internal/service"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	win       fyne.Window
	cfg       *model.AppConfig
	mover     *service.MouseMover
	interval  *widget.Entry
	toggleBtn *widget.Button
	v4App     fyne.App
}

func RunGUI() {
	cfg := config.LoadConfig()

	myApp := app.New()
	myApp.Settings().SetTheme(theme.DefaultTheme())

	window := myApp.NewWindow(cfg.GUI.Title)
	window.Resize(fyne.NewSize(float32(cfg.GUI.Width), float32(cfg.GUI.Height)))
	window.SetFixedSize(true)

	mover := service.NewMouseMover()

	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(fmt.Sprintf("%d", cfg.Mover.Interval))
	intervalEntry.SetPlaceHolder("单位:秒")
	intervalEntry.Validator = validation.NewRegexp(`^\d+$`, "请输入整数")

	toggleBtn := widget.NewButton("启动", nil)
	toggleBtn.Importance = widget.HighImportance

	g := &GUI{
		win:       window,
		cfg:       cfg,
		mover:     mover,
		interval:  intervalEntry,
		toggleBtn: toggleBtn,
		v4App:     myApp,
	}

	toggleBtn.OnTapped = g.onToggle

	content := container.NewVBox(
		container.NewPadded(
			container.NewBorder(nil, nil, widget.NewLabelWithStyle("时间间隔", fyne.TextAlignLeading, fyne.TextStyle{}), nil, intervalEntry),
		),
		container.NewPadded(
			container.NewBorder(nil, nil, nil, nil, toggleBtn),
		),
	)

	window.SetContent(content)
	window.SetOnClosed(func() {
		mover.Stop(func(msg string) {})
	})

	window.ShowAndRun()
}

func (g *GUI) onToggle() {
	if g.mover.IsRunning() {
		g.onStop()
	} else {
		g.onStart()
	}
}

func (g *GUI) onStart() {
	text := g.interval.Text
	interval, err := parseInt(text)
	if err != nil || interval < 1 {
		return
	}

	g.cfg.Mover.Interval = interval

	err = g.mover.Start(interval, func(msg string) {})
	if err != nil {
		return
	}

	g.toggleBtn.SetText("停止")
	g.toggleBtn.Importance = widget.DangerImportance
}

func (g *GUI) onStop() {
	g.mover.Stop(func(msg string) {})
	g.toggleBtn.SetText("启动")
	g.toggleBtn.Importance = widget.HighImportance
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
