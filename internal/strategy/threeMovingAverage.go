package strategy

import (
	"coinquant/internal/data"
	"coinquant/internal/strategy/module"
	"coinquant/pkg/upbit"
	"coinquant/pkg/upbit/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ThreeMovingAverage struct {
	first    int
	second   int
	third    int
	prev     map[string][]float64
	current  map[string][]float64
	position map[string]float64
	c        *upbit.Client
	db       *gorm.DB

	balance int
}

func InitStrategy(c *upbit.Client, db *gorm.DB, first, second, third int) *ThreeMovingAverage {
	prev := make(map[string][]float64)
	current := make(map[string][]float64)
	position := make(map[string]float64)

	prev["KRW-BTC"] = []float64{0, 0, 0}
	current["KRW-BTC"] = []float64{0, 0, 0}
	position["KRW-BTC"] = 0
	balance := 100000000
	return &ThreeMovingAverage{
		c:        c,
		db:       db,
		first:    first,
		second:   second,
		third:    third,
		prev:     prev,
		current:  current,
		position: position,
		balance:  balance,
	}
}

func (s *ThreeMovingAverage) Initialize() {

}

func (s *ThreeMovingAverage) TesterInitialize(period time.Duration) {
	today := time.Now()
	startTime := today.Add(-period)
	_, err := data.GetHistory(s.c, s.db, "KRW-BTC", model.Days, startTime, today)
	if err != nil {
		panic(err)
		return
	}
}

func (s *ThreeMovingAverage) Check(timeToCheck time.Time) {
	s.prev = make(map[string][]float64)
	for key, value := range s.current {
		copiedValue := make([]float64, len(value))
		copy(copiedValue, value)
		s.prev[key] = copiedValue
	}
	s.current = make(map[string][]float64)

	//timeToCheck := time.Now()

	//check positions
	for market, position := range s.position {
		first, _ := module.GetMovingAverage(s.c, s.db, market, model.Days, s.first, timeToCheck)
		second, _ := module.GetMovingAverage(s.c, s.db, market, model.Days, s.second, timeToCheck)
		third, _ := module.GetMovingAverage(s.c, s.db, market, model.Days, s.third, timeToCheck)
		s.current[market] = []float64{first, second, third}

		prev := s.prev[market]
		current := s.current[market]

		buy := buyCondition(position, prev, current)
		sell := sellCondition(position, prev, current)

		tobe := buy + sell
		if position != tobe {

			if buy > 0 && sell < 0 {
				fmt.Println("--------------")
				fmt.Println(timeToCheck.Format("2006-01-02T15:04:05"))
				fmt.Printf("%0.0f,%0.0f,%0.0f,\n%0.0f,%0.0f,%0.0f,\n", prev[0], prev[1], prev[2], current[0], current[1], current[2])
				fmt.Println(position, buy, sell)
			}
			if s.position[market]+tobe > 0 {
				//fmt.Printf("%s position : %f \n", market, position)
				s.position[market] = tobe
			}

			s.position[market] = tobe
		}

	}
}

func buyCondition(position float64, prev, current []float64) float64 {

	if position <= 0.0 && (prev[2] >= prev[1] && prev[1] >= prev[0]) && current[0] > current[1] {
		return 0.4
	}

	if position <= 0.4 && current[0] > current[2] {
		return 0.7
	}

	if position <= 0.7 && current[1] > current[2] {
		return 1.0
	}

	return 0.0
}

func sellCondition(position float64, prev, current []float64) float64 {
	if position <= 0.0 {
		return 0.0
	}
	if (prev[0] >= prev[1] && prev[1] >= prev[2]) && current[1] > current[0] {
		return -0.5
	} else if current[2] > current[0] {
		return -1.0
	}
	//
	//if current[1] < current[0] {
	//	return -1.0
	//}

	//if current[0] < current[1] {
	//	return -0.5
	//} else if current[0] < current[2] {
	//	return -1.0
	//}
	return 0.0
}
