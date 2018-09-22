package public

import (
	"encoding/json"

	"github.com/pkg/errors"

	"io/ioutil"
	"net/http"
	"time"
)

type Book struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type Orderbook struct {
	MidPrice float64 `json:"mid_price"`
	Bids     []Book  `json:"bids"`
	Asks     []Book  `json:"asks"`
}

const orderbook = "/v1/getboard"

func (p *Orderbook) Get() error {
	url := base + orderbook

	c := &http.Client{
		Timeout: 1000 * time.Millisecond,
	}

	res, err := c.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("error when gets orderbook")
	}

	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "read orderbook data from url: %s", url)
		}
		return errors.Wrapf(err, "unmarshal orderbook data: %s", string(data))
	}

	return nil
}
