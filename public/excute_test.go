package public

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

var E Executes

func TestGetExecute(t *testing.T) {
	var p = new(Executes)

	for i := 0; i < 2; i++ {
		var start = time.Now()
		func() {
			if err := p.Get(); err != nil {
				t.Error(err)
			}

			l := len(p.Execute) - 1
			E.LastID = p.Execute[l].ID
			E.Execute = append(E.Execute, p.Execute...)
		}()
		var ( // 1requestを500msに調整
			d   = 500 * time.Millisecond
			now = time.Now()
			sub = now.Sub(start)
		)
		fmt.Printf("%v - %v = %v\n", d, sub, d-sub)
		time.Sleep(d - sub)
	}

	sort.Sort(E)
	for _, v := range E.Execute {
		fmt.Printf("%+v\n", v.ID)
	}
	fmt.Printf("%+v\n", len(E.Execute))

}
