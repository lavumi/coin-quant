package data

import (
	"coinquant/pkg/upbit"
	"coinquant/pkg/upbit/model"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type missingInterval struct {
	end   time.Time
	count int
}

// todo fix 오늘 날짜가 포함된 요청이 있을 경우 문제가 있음
// 오늘의 데이터는 day chart 가 없으므로 요청해도 안온다.
// 그런데 calculateMissingIntervals 에서 오늘 날짜도 db에 없으므로
// 오늘 날짜에 대한 요청을 무의미하게 한번 더 하는 오류가 있음
func GetHistory(c *upbit.Client, db *gorm.DB, market string, interval model.CandleType, start time.Time, end time.Time) ([]model.Candle, error) {
	var existingCandles []model.Candle
	if interval == model.Days {
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		panic("only day intervals are supported")
	}

	err := db.Where("market = ? AND candle_date_time_utc BETWEEN ? AND ?", market, start, end.AddDate(0, 0, 1)).Find(&existingCandles).Error
	if err != nil {
		return nil, err
	}

	missingIntervals := calculateMissingIntervals(existingCandles, start, end)
	for _, missing := range missingIntervals {
		if err := storeDayChart(c, db, market, interval, missing.count, missing.end); err != nil {
			return nil, err
		}
	}
	db.Where("market = ? AND candle_date_time_utc BETWEEN ? AND ?", market, start, end).Find(&existingCandles)
	return existingCandles, nil
}

func calculateMissingIntervals(existingCandles []model.Candle, start, end time.Time) []missingInterval {
	var missingIntervals []missingInterval
	existingMap := make(map[string]bool)

	for _, candle := range existingCandles {
		existingMap[candle.CandleDateTimeUtc] = true
	}

	currentDate := start
	var currentInterval *missingInterval

	for currentDate.Before(end) || currentDate.Equal(end) {
		dateStr := currentDate.Format("2006-01-02T15:04:05")
		if !existingMap[dateStr] {
			fmt.Println("not exist ", dateStr)
			if currentInterval == nil {
				currentInterval = &missingInterval{end: currentDate, count: 1}
			} else {
				currentInterval.count++
				currentInterval.end = currentDate
			}
		} else {
			if currentInterval != nil {
				missingIntervals = append(missingIntervals, *currentInterval)
				currentInterval = nil
			}
		}
		currentDate = currentDate.Add(24 * time.Hour)
	}
	if currentInterval != nil {
		missingIntervals = append(missingIntervals, *currentInterval)
	}

	return missingIntervals
}

func storeDayChart(c *upbit.Client, db *gorm.DB, market string, interval model.CandleType, num int, end time.Time) error {
	const batchSize = 200
	days := num
	end = end.Add(time.Minute)
	for days > 0 {
		currentBatchSize := days
		if currentBatchSize > batchSize {
			currentBatchSize = batchSize
		}

		batchEnd := end
		batchStart := end.Add(-time.Duration(currentBatchSize*24) * time.Hour)
		candles, err := c.GetCandleChart(market, interval, batchEnd, currentBatchSize)
		if err != nil {
			return err
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "market"}, {Name: "candle_date_time_utc"}},
			UpdateAll: true,
		}).Create(&candles).Error; err != nil {
			return err
		}

		// Prepare for the next iteration
		end = batchStart
		days -= currentBatchSize
	}

	return nil
}
