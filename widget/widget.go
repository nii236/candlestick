package widget

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Sample is data for a single candle
type Sample struct {
	Time  time.Time
	Min   float64
	Max   float64
	Begin float64
	End   float64
}

// Candlestick is the main func used to create the widget
func Candlestick(getSamples func() []*Sample) *tview.Box {
	b := tview.NewBox()
	b.SetBorder(true)
	b.SetTitle("Candlestick")
	b.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		samples := getSamples()
		candles := []*Candle{}
		for _, s := range samples {
			candles = append(candles, NewCandle(s.Min, s.Max, s.Begin, s.End))
		}
		cc := NewCollection(candles, height)
		chartHeight := height - 2
		candleCount := width - 15 - 5
		if candleCount > len(cc.Data) {
			candleCount = len(cc.Data)
		}
		cc.Data = cc.Data[len(cc.Data)-candleCount : len(cc.Data)-1]
		for i, candle := range cc.Data {
			for cy := 2; cy < y+chartHeight; cy++ {
				ch := cc.renderCandleAt(candle, height-cy)
				rns := []rune(ch)
				r := rune(' ')
				if len(rns) > 0 {
					r = rns[len(rns)-1]
				}
				colour := tcell.ColorRed
				if candle.priceMove() == UpMove {
					colour = tcell.ColorGreen
				}
				if string(r) == SymbolNothing {
					colour = tcell.ColorWhite
				}
				screen.SetContent(i+x+15, cy, r, nil, tcell.StyleDefault.Foreground(colour))
			}
		}

		for cy := 0; cy < y+height-1; cy++ {
			tview.Print(screen, cc.renderAxesAt(height-cy), 0, cy, 15, tview.AlignRight, tcell.ColorYellow)
		}
		return b.GetInnerRect()
	})
	return b
}
