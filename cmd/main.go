package main

import (
	"coinquant/internal/strategy"
	"coinquant/internal/util"
	"coinquant/pkg/upbit"
	"coinquant/pkg/upbit/model"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	db, err := util.InitializeSqlite()
	if err != nil {
		log.Fatalf("fail to initialize db : %v", err)
	}

	client, err := upbit.MakeClient()
	if err != nil {
		log.Fatalf("fail to make upbit client : %v", err)
	}

	//err = strategy.InitializeCandleChart(client, db, model.Days, 50)
	//if err != nil {
	//	log.Fatalf("fail to make candle chart : %v", err)
	//}

	//res, err := util.GetAllCandles(db)

	//startTime := time.Date(2019, time.March, 1, 0, 0, 0, 0, time.UTC)

	//res, err := data.GetHistory(client, db, "KRW-BTC", model.Days, startTime, time.Now().UTC())
	//fmt.Println("=========================")
	//fmt.Println("Candles: ", len(res))
	//fmt.Println("=========================")

	movFive, _ := strategy.GetMovingAverage(client, db, "KRW-BTC", model.Days, 5, time.Now().UTC())
	//movTw, _ := strategy.GetMovingAverage(client, db, "KRW-BTC", model.Days, 20, time.Now().UTC())
	//movFty, _ := strategy.GetMovingAverage(client, db, "KRW-BTC", model.Days, 50, time.Now().UTC())

	fmt.Println("=========================")
	fmt.Printf("5 평균: %.0f\n", movFive)
	//fmt.Printf("20 평균: %.0f\n", movTw)
	//fmt.Printf("50 평균: %.0f\n", movFty)
	fmt.Println("=========================")

}
