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
	resp, err := client.Get("https://www.bitstamp.net/api/v2/ticker/btcusd/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Декодируем JSON.
	var tickerData models.TickerData
	if err := json.NewDecoder(resp.Body).Decode(&tickerData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	timestampStr := strconv.FormatInt(tickerData.Timestamp, 10)
	sideStr := strconv.Itoa(tickerData.Side)
	openStr := strconv.FormatFloat(tickerData.Open, 'f', -1, 64)
	highStr := strconv.FormatFloat(tickerData.High, 'f', -1, 64)
	lowStr := strconv.FormatFloat(tickerData.Low, 'f', -1, 64)
	lastStr := strconv.FormatFloat(tickerData.Last, 'f', -1, 64)
	volumeStr := strconv.FormatFloat(tickerData.Volume, 'f', -1, 64)
	vwapStr := strconv.FormatFloat(tickerData.VWAP, 'f', -1, 64)
	bidStr := strconv.FormatFloat(tickerData.Bid, 'f', -1, 64)
	askStr := strconv.FormatFloat(tickerData.Ask, 'f', -1, 64)
	open24Str := strconv.FormatFloat(tickerData.Open24, 'f', -1, 64)
	change24Str := strconv.FormatFloat(tickerData.PercentChange24, 'f', -1, 64)

	// Создаем новую запись.
	ticker := migrations.TradeRow{
		Timestamp:       timestampStr,
		Open:            openStr,
		High:            highStr,
		Low:             lowStr,
		Last:            lastStr,
		Volume:          volumeStr,
		VWAP:            vwapStr,
		Bid:             bidStr,
		Ask:             askStr,
		Side:            sideStr,
		Open24:          open24Str,
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