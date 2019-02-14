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
	fmt.Println(string('╽'))
	b := candlestick(SeedData())
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

func candlestick(cc CandleCollection) *tview.Box {
	b := tview.NewBox()

	// maxValue := 100
	// minValue := -100

	b.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for i, candle := range cc.data {
			ch := cc.renderCandleAt(candle, 10, true)

			screen.SetContent(i, 10, []rune(ch)[0], nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
		return b.GetInnerRect()
	})
	return b
}
