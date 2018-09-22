package public

import (
	"encoding/json"

	"github.com/pkg/errors"

	"io/ioutil"
	"net/http"
	"time"
)

type Ticker struct {
	Timestamp       BitflyerTime `json:"timestamp"`
	TickID          int          `json:"tick_id"`
	BestBid         float64      `json:"best_bid"`
	BestAsk         float64      `json:"best_ask"`
	BestBidSize     float64      `json:"best_bid_size"`
	BestAskSize     float64      `json:"best_ask_size"`
	TotalBidDepth   float64      `json:"total_bid_depth"`
	TotalAskDepth   float64      `json:"total_ask_depth"`
	LTP             float64      `json:"ltp"`
	Volume          float64      `json:"volume"`
	VolumeByProduct float64      `json:"volume_by_product"`
}

const ticker = "/v1/getticker"

func (p *Ticker) Get() error {
	url := base + ticker

	c := &http.Client{
		Timeout: 1000 * time.Millisecond,
	}

	res, err := c.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("error when gets ticker")
	}

	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "read ticker data from url: %s", url)
		}
		return errors.Wrapf(err, "unmarshal ticker data: %s", string(data))
	}

	return nil
}
