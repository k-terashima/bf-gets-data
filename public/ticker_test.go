package public

import (
	"fmt"
	"testing"
	"time"
)

func TestGetTicker(t *testing.T) {

	var p = new(Ticker)

	for i := 0; i < 10; i++ {
		var start = time.Now()
		func() {
			if err := p.Get(); err != nil {
				t.Error(err)
			}
		}()
		var ( // 1requestを500msに調整
			d   = 500 * time.Millisecond
			now = time.Now()
			sub = now.Sub(start)
		)
		fmt.Printf("%v - %v = %v\n", d, sub, d-sub)
		time.Sleep(d - sub)
	}
}
