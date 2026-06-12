package utils

import (
	"testing"
)

func TestBase(t *testing.T) {
	for i := 0; i < 3; i++ {
		defer func(){ println(i) } ()
	}
}
