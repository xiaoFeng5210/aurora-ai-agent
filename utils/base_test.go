package utils

import (
	"fmt"
	"testing"
)

func TestBase(t *testing.T) {
	for i := 0; i < 3; i++ {
		defer func() { println(i) }()
	}
}

func JaccardBySorted[T comparable](a, b []float64) float64 {
	if len(a) <= 0 || len(b) <= 0 {
		return 0.0
	}

	intersection := 0
	for i, j := 0, 0; i < len(a) && j < len(b); {
		if a[i] == b[j] {
			intersection += 1
			i += 1
			j += 1
		} else if a[i] < b[j] {
			i += 1
		} else if a[i] > b[j] {
			j += 1
		}
	}

	return float64(intersection) / float64(len(a)+len(b)-intersection)

}

func TestJaccardSimilarity(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	b := []float64{3.0, 4.0, 5.0, 6.0, 7.0}

	similarity := JaccardBySorted[float64](a, b)

	fmt.Println(similarity)

}
