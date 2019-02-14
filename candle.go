package main

import "math"

const SYMBOL_STICK = "│"

const SYMBOL_CANDLE = "┃"
const SYMBOL_HALF_TOP = "╽"
const SYMBOL_HALF_BOTTOM = "╿"
const SYMBOL_HALF_CANDLE_TOP = "╻"
const SYMBOL_HALF_CANDLE_BOTTOM = "╹"
const SYMBOL_HALF_STICK_TOP = "╷"
const SYMBOL_HALF_STICK_BOTTOM = "╵"
const SYMBOL_NOTHING = " "

const COLOR_NEUTRAL = "WHITE"
const COLOR_POSITIVE = "GREEN"
const COLOR_NEGATIVE = "RED"

const UP_MOVE = 1
const DOWN_MOVE = -1

type Candle struct {
	min_value   float64
	max_value   float64
	begin_value float64
	end_value   float64
}

func NewCandle(min_value float64, max_value float64, begin_value float64, end_value float64) *Candle {
	return &Candle{
		min_value:   min_value,
		max_value:   max_value,
		begin_value: begin_value,
		end_value:   end_value,
	}
}

func (c *Candle) topStick() float64 {
	return c.max_value
}

func (c *Candle) bottomStick() float64 {
	return c.min_value
}

func (c *Candle) topCandle() float64 {
	return math.Max(c.begin_value, c.end_value)
}

func (c *Candle) bottomCandle() float64 {
	return math.Min(c.begin_value, c.end_value)
}

func (c *Candle) priceMove() int {
	if c.begin_value > c.end_value {
		return DOWN_MOVE
	}
	return UP_MOVE
}

type CandleCollection struct {
	height         int
	data           []*Candle
	globalMinValue float64
	globalMaxValue float64
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
		height:         height,
		data:           candles,
		globalMinValue: globalMinValue,
		globalMaxValue: globalMaxValue,
	}

	return cc
}
func (cc *CandleCollection) toHeightUnits(x float64) float64 {
	return (x - cc.globalMinValue) / (cc.globalMaxValue - cc.globalMinValue) * float64(cc.height)
}

func (cc *CandleCollection) candleColor(candle *Candle) string {

	if candle.priceMove() == UP_MOVE {
		return COLOR_POSITIVE
	}

	return COLOR_NEGATIVE
}
func (cc *CandleCollection) renderCandleAt(candle *Candle, heightUnit int) string {
	heightUnit64 := float64(heightUnit)

	scaledTopStick := cc.toHeightUnits(candle.topStick())
	scaledTopCandle := cc.toHeightUnits(candle.topCandle())
	scaledBottomStick := cc.toHeightUnits(candle.bottomStick())
	scaledBottomCandle := cc.toHeightUnits(candle.bottomCandle())

	if math.Ceil(scaledTopStick) >= heightUnit64 && heightUnit64 >= math.Floor(scaledTopCandle) {
		if scaledTopCandle-heightUnit64 > 0.75 {
			return SYMBOL_CANDLE
		} else if (scaledTopCandle - heightUnit64) > 0.25 {
			if (scaledTopStick - heightUnit64) > 0.75 {
				return SYMBOL_HALF_TOP
			}
			return SYMBOL_HALF_CANDLE_TOP
		} else {
			if (scaledTopStick - heightUnit64) > 0.75 {
				return SYMBOL_STICK
			} else if (scaledTopStick - heightUnit64) > 0.25 {
				return SYMBOL_HALF_STICK_TOP
			} else {
				return SYMBOL_NOTHING
			}
		}
	} else if math.Floor(scaledTopCandle) >= heightUnit64 && heightUnit64 >= math.Ceil(scaledBottomCandle) {
		return SYMBOL_CANDLE
	} else if math.Ceil(scaledBottomCandle) >= heightUnit64 && heightUnit64 >= math.Floor(scaledBottomStick) {
		if (scaledBottomCandle - heightUnit64) < 0.25 {
			return SYMBOL_CANDLE
		} else if (scaledBottomCandle - heightUnit64) < 0.75 {
			if (scaledBottomStick - heightUnit64) < 0.25 {
				return SYMBOL_HALF_BOTTOM
			}
			return SYMBOL_HALF_CANDLE_BOTTOM
		} else {
			if (scaledBottomStick - heightUnit64) < 0.25 {
				return SYMBOL_STICK
			} else if (scaledBottomStick - heightUnit64) < 0.75 {
				return SYMBOL_HALF_STICK_BOTTOM
			} else {
				return SYMBOL_NOTHING
			}
		}
	}
	return SYMBOL_NOTHING
}

// func (cc *CandleCollection) draw(colorize bool) string {
// 	output_str = "\n"

// 	return output_str

// }

// for y in reversed(range(0, self._height)):
//     if y % 4 == 0:
//         output_str += (Style.RESET_ALL if colorize else "") + "{:8.2f} ".format(
//             self._global_min_value + (y * (self._global_max_value - self._global_min_value) / self._height))
//     else:
//         output_str += "         "
//     for c in self._data:
//         output_str += self._render_candle_at(c, y, colorize)
//     output_str += "\n" + (Style.RESET_ALL if colorize else "")
