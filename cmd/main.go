package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nii236/candlestick/widget"

	"github.com/gocarina/gocsv"
	"github.com/rivo/tview"
)

func main() {
	b := widget.Candlestick(getSamples)
	err := tview.NewApplication().SetRoot(b, true).Run()
	if err != nil {
		fmt.Println(err)
	}
}

func loadCSV() []*CSVRecord {
	ohlcFile, err := os.OpenFile("./finance-charts-apple.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer ohlcFile.Close()

	records := []*CSVRecord{}

	if err := gocsv.UnmarshalFile(ohlcFile, &records); err != nil {
		panic(err)
	}
	return records
}

// CSVRecord is used to marshal CSV records into a struct
type CSVRecord struct {
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

func getSamples() []*widget.Sample {
	records := loadCSV()
	samples := []*widget.Sample{}
	for _, rec := range records {
		layout := "2006-01-02"
		t, err := time.Parse(layout, rec.Date)
		if err != nil {
			panic(err)
		}
		samples = append(samples, &widget.Sample{
			Time:  t,
			Min:   rec.Low,
			Max:   rec.High,
			Begin: rec.Open,
			End:   rec.Close,
		})
	}

	return samples
	// return widget.SeedData()
}
