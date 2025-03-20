package trade

// calculateMACD вычисляет линии MACD, сигнальную линию и гистограмму.
// Для этого используются два EMA (например, 12 и 26 периодов) и EMA сигнальной линии (обычно 9 периодов).
func calculateMACD(prices []float64, shortPeriod, longPeriod, signalPeriod int) (macdLine, signalLine, histogram []float64) {
	emaShort := calculateEMA(prices, shortPeriod)
	emaLong := calculateEMA(prices, longPeriod)
	macdLine = make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		macdLine[i] = emaShort[i] - emaLong[i]
	}
	signalLine = calculateEMA(macdLine, signalPeriod)
	histogram = make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		histogram[i] = macdLine[i] - signalLine[i]
	}
	return
}
