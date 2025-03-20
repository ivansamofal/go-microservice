package trade

// calculateRSI рассчитывает индекс относительной силы для ряда цен с заданным периодом.
func calculateRSI(prices []float64, period int) []float64 {
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