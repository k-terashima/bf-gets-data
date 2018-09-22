package public

import (
	"net/http"
	"strings"
	"time"
)

const (
	base        = "https://api.bitflyer.com"
	productCode = "FX_BTC_JPY"
	count       = "500"
)

type Public interface {
	Get(*http.Client)
}

// Parse Bitflyer's time
type BitflyerTime struct {
	time.Time
}

const bitflyerTimeLayout = "2006-01-02T15:04:05.99999"

// changes bitflyerTime to time.Time
func (p *BitflyerTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "Z\"")
	p.Time, err = time.Parse(bitflyerTimeLayout, s)
	return err
}
