package module

import (
	"coinquant/internal/data"
	"coinquant/pkg/upbit"
	"coinquant/pkg/upbit/model"
	"gorm.io/gorm"
	"time"
)

func GetMovingAverage(c *upbit.Client, db *gorm.DB, market string, interval model.CandleType, limit int, end time.Time) (float64, error) {
	start := end.Add(-time.Duration(limit*24) * time.Hour)
	//get day chart
	candles, err := data.GetHistory(c, db, market, interval, start, end)
	if err != nil {
		return 0.0, err
	}

	//get current minutes data
	//chart, err := c.GetCandleChart(market, model.Minutes, end, 1)
	//if err != nil {
	//	return 0., err
	//}
	//candles = append(candles, chart[0])

	var sum float64
	for _, candle := range candles {
		sum += candle.TradePrice
	}
	return sum / float64(len(candles)), nil
}
