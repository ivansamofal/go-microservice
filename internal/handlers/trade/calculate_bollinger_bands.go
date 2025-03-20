package trade

import (
	"math"
)
// calculateBollingerBands рассчитывает линии скользящей средней и верхнюю/нижнюю полосы по заданному периоду и множителю стандартного отклонения.
func calculateBollingerBands(prices []float64, period int, k float64) (sma, upper, lower []float64) {
	sma = make([]float64, len(prices))
	upper = make([]float64, len(prices))
	lower = make([]float64, len(prices))
	for i := period - 1; i < len(prices); i++ {
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += prices[j]
		}
		avg := sum / float64(period)
		sma[i] = avg
		// Вычисляем стандартное отклонение.
		variance := 0.0
		for j := i - period + 1; j <= i; j++ {
			variance += math.Pow(prices[j]-avg, 2)
		}
		variance /= float64(period)
		stddev := math.Sqrt(variance)
		upper[i] = avg + k*stddev
		lower[i] = avg - k*stddev
	}
	return
}