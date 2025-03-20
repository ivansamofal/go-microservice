package trade

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"
	"crypto/tls"
	"go_microservice/internal/models"
	"go_microservice/internal/migrations"
	"go_microservice/internal/repository"
	"go_microservice/internal/logger"
	"github.com/gin-gonic/gin"
)

func SaveTradeData(c *gin.Context) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	// Выполняем HTTP GET запрос.
	resp, err := client.Get("https://www.binance.com/api/v3/ticker/24hr?symbol=BTCUSDT")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var tickerData models.BinanceTickerData
	if err := json.NewDecoder(resp.Body).Decode(&tickerData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	timestampSec := tickerData.CloseTime / 1000
	timestampStr := strconv.FormatInt(timestampSec, 10)

	//timestampStr := strconv.FormatInt(tickerData.Timestamp, 10)
	//sideStr := strconv.Itoa(tickerData.Side)
	openStr := strconv.FormatFloat(tickerData.OpenPrice, 'f', -1, 64)
	highStr := strconv.FormatFloat(tickerData.HighPrice, 'f', -1, 64)
	lowStr := strconv.FormatFloat(tickerData.LowPrice, 'f', -1, 64)
	lastStr := strconv.FormatFloat(tickerData.LastPrice, 'f', -1, 64)
	volumeStr := strconv.FormatFloat(tickerData.Volume, 'f', -1, 64)
	vwapStr := strconv.FormatFloat(tickerData.WeightedAvgPrice, 'f', -1, 64)
	bidStr := strconv.FormatFloat(tickerData.BidPrice, 'f', -1, 64)
	askStr := strconv.FormatFloat(tickerData.AskPrice, 'f', -1, 64)
	change24Str := strconv.FormatFloat(tickerData.PriceChangePercent, 'f', -1, 64)

	ticker := migrations.BinanceTicker {
		Timestamp:       timestampStr,
		Open:            openStr,
		High:            highStr,
		Low:             lowStr,
		Last:            lastStr,
		Volume:          volumeStr,
		VWAP:            vwapStr,
		Bid:             bidStr,
		Ask:             askStr,
		PercentChange24: change24Str,
		CreatedAt:       time.Now(),
	}

	err = repository.CreateTradeRow(&ticker)

	if err != nil {
		logger.Log.WithError(err).Error("Error during create trade row")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//// Сохраняем запись в БД.
	//if err := db.Create(&ticker).Error; err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"message": "Данные успешно сохранены",
		"data":    ticker,
	})
}