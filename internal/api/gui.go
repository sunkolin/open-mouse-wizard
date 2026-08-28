package api

import (
	"bling/internal/config"
	"bling/internal/model"
	"bling/internal/service"
	"fmt"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type GUI struct {
	win   *app.Window
	theme *material.Theme
	cfg   *model.AppConfig
	mover *service.MouseMover

	intervalInput widget.Editor
	startBtn      widget.Clickable
	stopBtn       widget.Clickable

	statusCh chan string
	status   string
	ops      op.Ops
}

func RunGUI() {
	cfg := config.LoadConfig()

	window := &app.Window{}
	window.Option(
		app.Title(cfg.GUI.Title),
		app.Size(unit.Dp(300), unit.Dp(160)),
		app.MinSize(unit.Dp(300), unit.Dp(160)),
	)

	mover := service.NewMouseMover()

	shaper := text.NewShaper(text.WithCollection(gofont.Collection()))
	theme := material.NewTheme()
	theme.Shaper = shaper

	g := &GUI{
		win:      window,
		theme:    theme,
		cfg:      cfg,
		mover:    mover,
		statusCh: make(chan string, 16),
	}

	g.intervalInput.SetText(fmt.Sprintf("%d", cfg.Mover.Interval))
	g.status = "就绪"

	go g.loop()
}

func (g *GUI) loop() {
	for {
		select {
		case msg := <-g.statusCh:
			g.status = msg
			g.win.Invalidate()
		default:
		}

		e := g.win.Event()
		switch e := e.(type) {
		case app.FrameEvent:
			gtx := layout.Context{
				Ops:    &g.ops,
				Source: e.Source,
				Metric: e.Metric,
				Now:    e.Now,
				Constraints: layout.Constraints{
					Max: e.Size,
				},
			}
			g.layout(gtx)
			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			g.mover.Stop(func(msg string) {})
			return
		}
	}
}

func (g *GUI) layout(gtx layout.Context) {
	paint.ColorOp{Color: g.theme.Palette.Bg}.Add(gtx.Ops)

	if g.startBtn.Clicked(gtx) {
		g.onStart()
	}
	if g.stopBtn.Clicked(gtx) {
		g.onStop()
	}

	running := g.mover.IsRunning()

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body1(g.theme, "时间间隔(秒)")
								return lbl.Layout(gtx)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Dimensions{}
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Caption(g.theme, g.status)
								if running {
									lbl.Color = color.NRGBA{R: 0, G: 128, B: 0, A: 255}
								}
								return lbl.Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return g.inputField(gtx)
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if running {
							gtx = gtx.Disabled()
						}
						btn := material.Button(g.theme, &g.startBtn, "启动")
						btn.CornerRadius = 8
						btn.Background = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
						return btn.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{}
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if !running {
							gtx = gtx.Disabled()
						}
						btn := material.Button(g.theme, &g.stopBtn, "停止")
						btn.CornerRadius = 8
						btn.Background = color.NRGBA{R: 200, G: 60, B: 60, A: 255}
						return btn.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (g *GUI) onStart() {
	interval, err := parseInt(g.intervalInput.Text())
	if err != nil || interval < 1 {
		g.status = "请输入≥1的整数"
		return
	}

	g.cfg.Mover.Interval = interval

	err = g.mover.Start(interval, func(msg string) {
		g.statusCh <- msg
	})
	if err != nil {
		g.status = fmt.Sprintf("启动失败: %v", err)
	}
}

func (g *GUI) onStop() {
	g.mover.Stop(func(msg string) {
		g.statusCh <- msg
	})
}

func (g *GUI) inputField(gtx layout.Context) layout.Dimensions {
	g.intervalInput.SingleLine = true

	borderColor := color.NRGBA{R: 180, G: 180, B: 180, A: 255}
	inputBg := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			rr := gtx.Dp(unit.Dp(4))
			defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, inputBg)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		func(gtx layout.Context) layout.Dimensions {
			return widget.Border{
				Color:        borderColor,
				Width:        unit.Dp(1),
				CornerRadius: unit.Dp(4),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					ed := material.Editor(g.theme, &g.intervalInput, "60")
					ed.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
					ed.HintColor = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
					return ed.Layout(gtx)
				})
			})
		},
	)
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
