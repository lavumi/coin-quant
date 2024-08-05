package strategy

import "testing"

func TestBuyCondition_CheckInverseArray(t *testing.T) {
	current := []float64{1.0, 2.0, 3.0}
	prev := []float64{3.0, 2.0, 3.0}
	res2 := buyCondition(prev, current)

	prev = []float64{1.0, 2.0, 1.0}
	res3 := buyCondition(prev, current)

	if res2 != 0 {
		t.Error("res2 should be 0", res2)
	}

	if res3 != 0 {
		t.Error("res3 should be 0 : ", res3)
	}
}

func TestBuyCondition_CheckFirstBuy(t *testing.T) {
	prev := []float64{1.0, 2.0, 3.0}
	current := []float64{3.0, 2.0, 3.0}
	condition := buyCondition(prev, current)
	if condition != 0.4 {
		t.Error("res1 should be 0.4 but", condition)
	}
}

func TestBuyCondition_CheckSecond(t *testing.T) {
	prev := []float64{3.0, 2.0, 3.0}
	current := []float64{4.0, 2.0, 3.0}
	condition := buyCondition(prev, current)
	if condition != 0.7 {
		t.Error("res1 should be 0.7 but", condition)
	}
}

func TestBuyCondition_CheckThird(t *testing.T) {
	prev := []float64{4.0, 2.0, 3.0}
	current := []float64{4.0, 3.5, 3.0}
	condition := buyCondition(prev, current)
	if condition != 1.0 {
		t.Error("res1 should  be 1.0 but", condition)
	}
}
