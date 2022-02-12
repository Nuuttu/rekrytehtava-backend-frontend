package main

import (
	"testing"
)

// Testejä ei oikee kerenny. Meni yliaikaa
func TestContainsId(t *testing.T) {
	i := 20
	ints := []int{10, 20, 30, 40, 50}

	if !containsId(ints, 20) {
		t.Errorf("contains(ints, 20) = %d; want true", i)
	}

}

func TestVerifyCoffee(t *testing.T) {
	n := "Name"
	w := 200
	p := 20.20
	r := 5

	cof := Coffee{
		Name:       &n,
		Weight:     &w,
		Price:      &p,
		RoastLevel: &r,
	}

	err := verifyCoffee(cof)
	if err == nil {
		t.Errorf("verify(cof) = %d; want %d", err, nil)
	}
}
