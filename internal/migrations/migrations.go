package migrations

import (
	"time"
)

type BinanceTicker struct {
	ID              uint      `gorm:"primaryKey"`
	Timestamp       string    `gorm:"column:timestamp;not null"`        // используем CloseTime/1000 как Unix timestamp
	Open            string    `gorm:"column:open;not null"`             // соответствует openPrice
	High            string    `gorm:"column:high;not null"`             // соответствует highPrice
	Low             string    `gorm:"column:low;not null"`              // соответствует lowPrice
	Last            string    `gorm:"column:last;not null"`             // соответствует lastPrice
	Volume          string    `gorm:"column:volume;not null"`           // соответствует volume
	VWAP            string    `gorm:"column:vwap;not null"`             // соответствует weightedAvgPrice
	Bid             string    `gorm:"column:bid;not null"`              // соответствует bidPrice
	Ask             string    `gorm:"column:ask;not null"`              // соответствует askPrice
	PercentChange24 string    `gorm:"column:percent_change_24;not null"`// соответствует priceChangePercent
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
}
