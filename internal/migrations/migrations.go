package migrations

import (
	"time"
)
type TradeRow struct {
	ID              uint      `gorm:"primaryKey"`
	Timestamp       string    `gorm:"column:timestamp"`
	Open            string    `gorm:"column:open"`
	High            string    `gorm:"column:high"`
	Low             string    `gorm:"column:low"`
	Last            string    `gorm:"column:last"`
	Volume          string    `gorm:"column:volume"`
	VWAP            string    `gorm:"column:vwap"`
	Bid             string    `gorm:"column:bid"`
	Ask             string    `gorm:"column:ask"`
	Side            string    `gorm:"column:side"`
	Open24          string    `gorm:"column:open_24"`
	PercentChange24 string    `gorm:"column:percent_change_24"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}
