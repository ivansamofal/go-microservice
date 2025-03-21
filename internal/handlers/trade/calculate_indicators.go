package trade

import (
	"go_microservice/internal/models"
)

// calculateEMA вычисляет экспоненциальную скользящую среднюю по ряду цен.
// Если длина ряда меньше периода, возвращается срез с нулями.
func CalculateEMA(prices []float64, period int) []float64 {
	ema := make([]float64, len(prices))
	if len(prices) < period {
		return ema
	}
	// Первое значение EMA равно простому скользящему среднему за период.
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)
	ema[period-1] = sma

	multiplier := 2.0 / (float64(period) + 1)
	// Вычисляем EMA для последующих значений.
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i]-ema[i-1])*multiplier + ema[i-1]
	}
	return ema
}

// calculateMACD вычисляет линии MACD, сигнальную линию и гистограмму.
// Для этого используются два EMA (например, 12 и 26 периодов) и EMA сигнальной линии (обычно 9 периодов).
func CalculateMACD(prices []float64, shortPeriod int, longPeriod int, signalPeriod int) (macdLine, signalLine, histogram []float64) {
	emaShort := CalculateEMA(prices, shortPeriod)
	emaLong := CalculateEMA(prices, longPeriod)
	macdLine = make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		macdLine[i] = emaShort[i] - emaLong[i]
	}
	signalLine = CalculateEMA(macdLine, signalPeriod)
	histogram = make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		histogram[i] = macdLine[i] - signalLine[i]
	}
	return
}

func CalculateVWAP(data []models.BinanceTickerData) float64 {
	totalPV, totalVolume := 0.0, 0.0
	for _, d := range data {
		totalPV += d.LastPrice * d.Volume
		totalVolume += d.Volume
	}
	if totalVolume == 0 {
		return 0
	}
	return totalPV / totalVolume
}

// calculateRSI рассчитывает индекс относительной силы для ряда цен с заданным периодом.
func CalculateRSI(prices []float64, period int) []float64 {
	rsi := make([]float64, len(prices))
	if len(prices) < period+1 {
		return rsi
	}
	gains, losses := 0.0, 0.0
	// Начальные средние приросты и потери.
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses += -change
		}
	}
	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		rsi[period] = 100
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}
	// Вычисление RSI для последующих периодов.
	for i := period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		gain, loss := 0.0, 0.0
		if change > 0 {
			gain = change
		} else {
			loss = -change
		}
		avgGain = ((avgGain * float64(period-1)) + gain) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + loss) / float64(period)
		var rs float64
		if avgLoss == 0 {
			rs = 0
		} else {
			rs = avgGain / avgLoss
		}
		rsi[i] = 100 - (100 / (1 + rs))
	}
	return rsi
}