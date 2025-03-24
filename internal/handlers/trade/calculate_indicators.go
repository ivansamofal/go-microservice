package trade

import (
	"go_microservice/internal/models"
	//"go_microservice/internal/logger"
	//"strconv"
	"math"
)
// fullEMA вычисляет полную серию значений EMA без обрезки.
func fullEMA(prices []float64, period int) []float64 {
	ema := make([]float64, len(prices))
	if len(prices) < period {
		// Если данных меньше периода, вернуть пустой срез (или можно вернуть ema, заполненный нулями).
		return ema
	}
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)
	ema[period-1] = sma

	multiplier := 2.0 / (float64(period) + 1)
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i]-ema[i-1])*multiplier + ema[i-1]
	}
	return ema
}

// CalculateEMA возвращает последние 30 элементов рассчитанной EMA (если их достаточно),
// иначе возвращает всю серию.
func CalculateEMA(prices []float64, period int) []float64 {
	ema := fullEMA(prices, period)
	if len(ema) >= 30 {
		return ema[len(ema)-30:]
	}
	return ema
}

// CalculateMACD вычисляет MACD, сигнальную линию и гистограмму.
// Для расчёта используются полные серии EMA, после чего результаты обрезаются до последних 30 элементов.
func CalculateMACD(prices []float64, shortPeriod, longPeriod, signalPeriod int) (macdLine, signalLine, histogram []float64) {
	emaShort := fullEMA(prices, shortPeriod)
	emaLong := fullEMA(prices, longPeriod)

	macdFull := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		macdFull[i] = emaShort[i] - emaLong[i]
	}

	signalFull := fullEMA(macdFull, signalPeriod)
	histogramFull := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		histogramFull[i] = macdFull[i] - signalFull[i]
	}

	if len(macdFull) >= 30 {
		macdLine = macdFull[len(macdFull)-30:]
		signalLine = signalFull[len(signalFull)-30:]
		histogram = histogramFull[len(histogramFull)-30:]
	} else {
		macdLine = macdFull
		signalLine = signalFull
		histogram = histogramFull
	}
	return
}

// CalculateVWAP рассчитывает VWAP на основе переданных данных.
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

// CalculateRSI рассчитывает RSI для ряда цен с заданным периодом.
// После расчёта возвращается срез из последних 30 значений (если данных достаточно).
func CalculateRSI(prices []float64, period int) []float64 {
	rsi := make([]float64, len(prices))
	if len(prices) < period+1 {
		return rsi
	}
	gains, losses := 0.0, 0.0
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
			rs = math.Inf(1)
		} else {
			rs = avgGain / avgLoss
		}
		rsi[i] = 100 - (100 / (1 + rs))
	}

	if len(rsi) >= 30 {
		return rsi[len(rsi)-30:]
	}
	return rsi
}