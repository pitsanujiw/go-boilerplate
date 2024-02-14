package memstore

import (
	"fmt"
	"testing"
)

func TestExpBackoff(t *testing.T) {

	exp := NewExponentialBackOff()

	for i := 0; i < 25; i++ {
		d := exp.NextBackOff()
		fmt.Printf("try: %v; next: %v\n", i, d)
	}

}
