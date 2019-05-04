package widget

import (
	"fmt"
	"math"
)

type Candle struct {
	MinValue   float64
	MaxValue   float64
	BeginValue float64
	EndValue   float64
}

func NewCandle(MinValue float64, MaxValue float64, BeginValue float64, EndValue float64) *Candle {
	return &Candle{
		MinValue:   MinValue,
		MaxValue:   MaxValue,
		BeginValue: BeginValue,
		EndValue:   EndValue,
	}
}

func (c *Candle) topStick() float64 {
	return c.MaxValue
}

func (c *Candle) bottomStick() float64 {
	return c.MinValue
}

func (c *Candle) topCandle() float64 {
	return math.Max(c.BeginValue, c.EndValue)
}

func (c *Candle) bottomCandle() float64 {
	return math.Min(c.BeginValue, c.EndValue)
}

func (c *Candle) priceMove() int {
	if c.BeginValue > c.EndValue {
		return DownMove
	}
	return UpMove
}

type CandleCollection struct {
	Height         int
	Data           []*Candle
	GlobalMinValue float64
	GlobalMaxValue float64
}

func NewCollection(candles []*Candle, height int) CandleCollection {
	globalMinValue := math.MaxFloat64
	for _, candle := range candles {
		if candle.bottomStick() < globalMinValue {
			globalMinValue = candle.bottomStick()
		}
	}

	globalMaxValue := float64(math.MinInt64)
	for _, candle := range candles {
		if candle.topStick() > globalMaxValue {
			globalMaxValue = candle.topStick()
		}
	}

	cc := CandleCollection{
		Height:         height,
		Data:           candles,
		GlobalMinValue: globalMinValue,
		GlobalMaxValue: globalMaxValue,
	}

	return cc
}
func (cc *CandleCollection) toHeightUnits(x float64) float64 {
	return (x - cc.GlobalMinValue) / (cc.GlobalMaxValue - cc.GlobalMinValue) * float64(cc.Height)
}

func (cc *CandleCollection) candleColor(candle *Candle) string {

	if candle.priceMove() == UpMove {
		return ColorPositive
	}

	return ColorNegative
}
func (cc *CandleCollection) renderCandleAt(candle *Candle, heightUnit int) string {
	heightUnit64 := float64(heightUnit)

	scaledTopStick := cc.toHeightUnits(candle.topStick())
	scaledTopCandle := cc.toHeightUnits(candle.topCandle())
	scaledBottomStick := cc.toHeightUnits(candle.bottomStick())
	scaledBottomCandle := cc.toHeightUnits(candle.bottomCandle())

	if math.Ceil(scaledTopStick) >= heightUnit64 && heightUnit64 >= math.Floor(scaledTopCandle) {
		if scaledTopCandle-heightUnit64 > 0.75 {
			return SymbolCandle
		} else if (scaledTopCandle - heightUnit64) > 0.25 {
			if (scaledTopStick - heightUnit64) > 0.75 {
				return SymbolHalfTop
			}
			return SymbolHalfCandleTop
		} else {
			if (scaledTopStick - heightUnit64) > 0.75 {
				return SymbolStick
			} else if (scaledTopStick - heightUnit64) > 0.25 {
				return SymbolHalfStickTop
			} else {
				return SymbolNothing
			}
		}
	} else if math.Floor(scaledTopCandle) >= heightUnit64 && heightUnit64 >= math.Ceil(scaledBottomCandle) {
		return SymbolCandle
	} else if math.Ceil(scaledBottomCandle) >= heightUnit64 && heightUnit64 >= math.Floor(scaledBottomStick) {
		if (scaledBottomCandle - heightUnit64) < 0.25 {
			return SymbolCandle
		} else if (scaledBottomCandle - heightUnit64) < 0.75 {
			if (scaledBottomStick - heightUnit64) < 0.25 {
				return SymbolHalfBottom
			}
			return SymbolHalfCandleBottom
		} else {
			if (scaledBottomStick - heightUnit64) < 0.25 {
				return SymbolStick
			} else if (scaledBottomStick - heightUnit64) < 0.75 {
				return SymbolHalfStickBottom
			} else {
				return SymbolNothing
			}
		}
	}
	return SymbolNothing
}

func (cc *CandleCollection) renderAxesAt(y int) string {
	if y%4 == 0 {
		return fmt.Sprintf("%.2f", cc.GlobalMinValue+(float64(y)*(cc.GlobalMaxValue-cc.GlobalMinValue)/float64(cc.Height)))
	}
	return ""
}

// for y in reversed(range(0, self._height)):
//     if y % 4 == 0:
//         output_str += (Style.RESET_ALL if colorize else "") + "{:8.2f} ".format(
//             self._global_min_value + (y * (self._global_max_value - self._global_min_value) / self._height))
//     else:
//         output_str += "         "
//     for c in self._data:
//         output_str += self._render_candle_at(c, y, colorize)
//     output_str += "\n" + (Style.RESET_ALL if colorize else "")
