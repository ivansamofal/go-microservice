package trade

// calculateEMA вычисляет экспоненциальную скользящую среднюю по ряду цен.
// Если длина ряда меньше периода, возвращается срез с нулями.
func calculateEMA(prices []float64, period int) []float64 {
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