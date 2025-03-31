package controllers

import (
	"go_microservice/internal/handlers/trade"
	"go_microservice/internal/repository"
	"go_microservice/internal/logger"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go_microservice/internal/models"
	"go_microservice/internal/migrations"
)

func Calculations(c *gin.Context) {
	x, err := strconv.Atoi(c.DefaultQuery("x", "1"))
	if err != nil || x <= 0 {
		x = 1
	}
	// Получаем период для EMA из query-параметра, по умолчанию 10
	period, err := strconv.Atoi(c.DefaultQuery("p", "10"))
	if err != nil || period <= 0 {
		period = 10
	}
	shortPeriod, err := strconv.Atoi(c.DefaultQuery("sp", "12"))
	if err != nil || shortPeriod <= 0 {
		shortPeriod = 12
	}
	longPeriod, err := strconv.Atoi(c.DefaultQuery("lp", "26"))
	if err != nil || longPeriod <= 0 {
		longPeriod = 26
	}
	signalPeriod, err := strconv.Atoi(c.DefaultQuery("sinp", "9"))
	if err != nil || signalPeriod <= 0 {
		signalPeriod = 9
	}

	var tickers []migrations.BinanceTicker
	err = repository.GetTradeRows(&tickers, longPeriod * x)
	if err != nil {
		logger.Log.WithError(err).Error("Error during retrieving trade rows")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(tickers) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No data"})
		return
	}

	// Подготовка данных для расчёта EMA и VWAP
	var prices []float64
	var tickerData []models.BinanceTickerData

	for _, t := range tickers {
		// Преобразуем поле Last из string в float64
		lastPrice, err := strconv.ParseFloat(t.Last, 64)
		if err != nil {
			continue // или обработать ошибку
		}
		prices = append(prices, lastPrice)

		vol, err := strconv.ParseFloat(t.Volume, 64)
		if err != nil {
			vol = 0
		}
		tData := models.BinanceTickerData{
			LastPrice: lastPrice,
			Volume:    vol,
		}
		tickerData = append(tickerData, tData)
	}

	// Вычисляем EMA
	ema := trade.CalculateEMA(prices, period)

	// Вычисляем VWAP
	vwap := trade.CalculateVWAP(tickerData)

	// Вычисляем MACD (macdLine, сигнальная линия и гистограмма)
	macdLine, signalLine, histogram := trade.CalculateMACD(prices, shortPeriod, longPeriod, signalPeriod)

	// Вычисляем RSI
	rsi := trade.CalculateRSI(prices, period)

	c.JSON(http.StatusOK, gin.H{
		"ema":        ema,
		"vwap":       vwap,
		"macdLine":   macdLine,
		"signalLine": signalLine,
		"histogram":  histogram,
		"rsi":        rsi,
	})
}
