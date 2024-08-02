package model

type CandleType string

const (
	Minutes CandleType = "minutes/1"
	Days    CandleType = "days"
	Weeks   CandleType = "weeks"
	Months  CandleType = "months"
)

type Candle struct {
	Market               string  `json:"market" gorm:"primaryKey"`
	CandleDateTimeUtc    string  `json:"candle_date_time_utc" gorm:"primaryKey"`
	CandleDateTimeKst    string  `json:"candle_date_time_kst"`
	OpeningPrice         float64 `json:"opening_price"`
	HighPrice            float64 `json:"high_price"`
	LowPrice             float64 `json:"low_price"`
	TradePrice           float64 `json:"trade_price"`
	Timestamp            int64   `json:"timestamp"`
	CandleAccTradePrice  float64 `json:"candle_acc_trade_price"`
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	Unit                 int     `json:"unit"`
}
