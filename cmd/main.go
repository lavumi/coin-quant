package main

import (
	"coinquant/internal/strategy"
	"coinquant/internal/util"
	"coinquant/pkg/upbit"
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

	//movFive, _ := module.GetMovingAverage(client, db, "KRW-BTC", model.Days, 5, time.Now().UTC())
	//movTw, _ := module.GetMovingAverage(client, db, "KRW-BTC", model.Days, 20, time.Now().UTC())
	//movFty, _ := module.GetMovingAverage(client, db, "KRW-BTC", model.Days, 50, time.Now().UTC())
	//
	//fmt.Println("=========================")
	//fmt.Printf("5 평균: %.0f\n", movFive)
	//fmt.Printf("20 평균: %.0f\n", movTw)
	//fmt.Printf("50 평균: %.0f\n", movFty)
	//fmt.Println("=========================")

	test := strategy.InitStrategy(client, db, 5, 20, 50)

	test.TesterInitialize(700 * 24 * time.Hour)
	pastDate := time.Now().AddDate(0, 0, -700)
	for i := 0; i < 700; i++ {
		test.Check(pastDate)
		pastDate = pastDate.AddDate(0, 0, 1)
		//time.Sleep(1000 * time.Millisecond)
	}
}
