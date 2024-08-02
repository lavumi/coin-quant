package data

import (
	"coinquant/pkg/upbit/model"
	"testing"
	"time"
)

func TestCalculateMissingIntervals(t *testing.T) {
	layout := "2006-01-02T15:04:05"
	start, _ := time.Parse(layout, "2023-01-01T00:00:00")
	end, _ := time.Parse(layout, "2023-01-10T00:00:00")

	existingCandles := []model.Candle{
		{CandleDateTimeUtc: "2023-01-01T00:00:00"},
		{CandleDateTimeUtc: "2023-01-03T00:00:00"},
		{CandleDateTimeUtc: "2023-01-05T00:00:00"},
	}

	expected := []missingInterval{
		{end: time.Date(2023, 01, 02, 0, 0, 0, 0, time.UTC), count: 1},
		{end: time.Date(2023, 01, 04, 0, 0, 0, 0, time.UTC), count: 1},
		{end: time.Date(2023, 01, 10, 0, 0, 0, 0, time.UTC), count: 5},
	}
	result := calculateMissingIntervals(existingCandles, start, end)
	if len(result) != len(expected) {
		t.Errorf("expected %d missing intervals, got %d", len(expected), len(result))
	}

	for i, interval := range result {
		if !interval.end.Equal(expected[i].end) || interval.count != expected[i].count {
			t.Errorf("expected interval %v, got %v", expected[i], interval)
		}
	}
}
