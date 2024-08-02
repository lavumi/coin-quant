package util

import (
	"coinquant/pkg/upbit/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func InitializeSqlite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("quant.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&model.Candle{})
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func GetAllCandles(db *gorm.DB) ([]model.Candle, error) {
	var candles []model.Candle
	result := db.Find(&candles)
	return candles, result.Error
}
