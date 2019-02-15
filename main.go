package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gocarina/gocsv"

	"github.com/rivo/tview"
)

type Record struct {
	Date      string  `csv:"Date"`
	Open      float64 `csv:"AAPL.Open"`
	High      float64 `csv:"AAPL.High"`
	Low       float64 `csv:"AAPL.Low"`
	Close     float64 `csv:"AAPL.Close"`
	Volume    float64 `csv:"AAPL.Volume"`
	Adjusted  float64 `csv:"AAPL.Adjusted"`
	dn        float64 `csv:"dn"`
	mavg      float64 `csv:"mavg"`
	up        float64 `csv:"up"`
	direction string  `csv:"direction"`
}

func main() {

	// for _, c := range data.data {
	// 	fmt.Print(data.renderCandleAt(c, 0, true))
	// }
	b := candlestick()
	err := tview.NewApplication().SetRoot(b, true).Run()
	if err != nil {

		fmt.Println(err)
	}
}

func load() []*Record {
	ohlcFile, err := os.OpenFile("./finance-charts-apple.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer ohlcFile.Close()

	records := []*Record{}

	if err := gocsv.UnmarshalFile(ohlcFile, &records); err != nil { // Load records from file
		panic(err)
	}
	for _, r := range records {
		fmt.Println("Hello", r.Adjusted)
	}
	return records
}

func candlestick() *tview.Box {
	b := tview.NewBox()

	// maxValue := 100
	// minValue := -100
	b.SetBorder(true)
	b.SetTitle("Candlestick pre-alpha")
	b.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		chartHeight := height - 2
		cc := SeedData(height)
		candleCount := width - 15 - 5
		if candleCount > len(cc.data) {
			candleCount = len(cc.data)
		}
		cc.data = cc.data[len(cc.data)-candleCount : len(cc.data)-1]
		for i, candle := range cc.data {
			for cy := 2; cy < y+chartHeight; cy++ {
				ch := cc.renderCandleAt(candle, height-cy)
				rns := []rune(ch)
				r := rune(' ')
				if len(rns) > 0 {
					r = rns[len(rns)-1]
				}
				colour := tcell.ColorRed
				if candle.priceMove() == UP_MOVE {
					colour = tcell.ColorGreen
				}
				if string(r) == SYMBOL_NOTHING {
					colour = tcell.ColorWhite
				}
				// if string(r) != SYMBOL_NOTHING {
				// 	fmt.Println(string(r))
				// }
				screen.SetContent(i+x+15, cy, r, nil, tcell.StyleDefault.Foreground(colour))
			}
			// fmt.Println(candle.end_value)
			// break
		}

		for cy := 0; cy < y+height-1; cy++ {
			tview.Print(screen, cc.renderAxesAt(height-cy), 0, cy, 15, tview.AlignRight, tcell.ColorYellow)
		}
		return b.GetInnerRect()
	})
	return b
}
