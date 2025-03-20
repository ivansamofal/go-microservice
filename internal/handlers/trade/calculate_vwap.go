package trade

import (
	"go_microservice/internal/models"
)

// calculateVWAP рассчитывает средневзвешенную по объёму цену (VWAP) за день.
func calculateVWAP(data []models.TickerData) float64 {
	totalPV, totalVolume := 0.0, 0.0
	for _, d := range data {
		totalPV += d.Last * d.Volume
		totalVolume += d.Volume
	}
	if totalVolume == 0 {
		return 0
	}
	return totalPV / totalVolume
}
