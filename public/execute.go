package public

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"

	"net/http"
	"sort"
	"strconv"
	"time"
)

type Executes struct {
	LastID  int
	Execute []Execute
}

type Execute struct {
	ID                         int          `json:"id"`
	Side                       string       `json:"side"`
	Price                      float64      `json:"price"`
	Size                       float64      `json:"size"`
	ExecDate                   BitflyerTime `json:"exec_date"`
	BuyChildOrderAcceptanceID  string       `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string       `json:"sell_child_order_acceptance_id"`
}

const execute = "/v1/getexecutions"

func (p Executes) Len() int {
	return len(p.Execute)
}

func (p Executes) Swap(i, j int) {
	p.Execute[i], p.Execute[j] = p.Execute[j], p.Execute[i]
}

func (p Executes) Less(i, j int) bool {
	return p.Execute[i].ID < p.Execute[j].ID
}

func (p *Executes) Get() error {
	var (
		url string
	)
	if p.LastID != 0 {
		s := strconv.Itoa(p.LastID)
		url = base + execute + "?product_code=" + productCode + "&count=" + count + "&after=" + s
	} else {
		url = base + execute + "?product_code=" + productCode + "&count=" + count
	}

	c := &http.Client{
		Timeout: 1000 * time.Millisecond,
	}

	res, err := c.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("error when gets executions: " + res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(&p.Execute); err != nil {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "read execute data from url: %s", url)
		}
		return errors.Wrapf(err, "unmarshal execute data: %s", string(data))
	}

	// sort id
	sort.Sort(*p)

	// save lastID
	l := len(p.Execute) - 1
	if l > 0 {
		p.LastID = p.Execute[l].ID
	}

	return nil
}
