package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	firebase "firebase.google.com/go"
	"github.com/k-terashima/bf-gets-data/public"
	"github.com/labstack/gommon/log"
	"google.golang.org/api/option"
)

type Data struct {
	Tickers    public.Ticker
	Orderbooks public.Orderbook
	Executes   []public.Execute

	CreatedAt time.Time
}

func main() {
	var (
		D         Data
		wg        sync.WaitGroup
		ticker    = new(public.Ticker)
		orderbook = new(public.Orderbook)
		execute   = new(public.Executes)

		ti []public.Ticker
		or []public.Orderbook
	)

	for i := 0; i < 501; i++ {
		var (
			start               = time.Now()
			fix   time.Duration = 500 * time.Millisecond
		)

		func(i int) {
			wg.Add(1)
			go func() {
				if err := ticker.Get(); err != nil {
					log.Error(err)
				}
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				if err := orderbook.Get(); err != nil {
					log.Error(err)
				}
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				if err := execute.Get(); err != nil {
					log.Error(err)
				}
				wg.Done()
			}()

			wg.Wait()

			D.Tickers = *ticker
			D.Orderbooks = *orderbook
			D.Executes = execute.Execute
			D.CreatedAt = time.Now()

			go func(i int, D Data) {
				ti = append(ti, D.Tickers)
				or = append(or, D.Orderbooks)
				if i%500 == 0 {
					if err := uploadStrage("ticker", ti); err != nil {
						log.Error(err)
						return
					}
					if err := uploadStrage("orderbook", or); err != nil {
						log.Error(err)
						return
					}
				}

				if err := uploadStrage("execute", D.Executes); err != nil {
					log.Error(err)
					return
				}
			}(i, D)

			// reset data
			D = Data{}
		}(i)

		end := time.Now()
		wait := fix - end.Sub(start)
		if wait < time.Duration(0) {
			wait = 0
		}
		time.Sleep(wait)
	}

}

func uploadStrage(where string, o interface{}) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		fmt.Println("database: ", end.Sub(start))
	}()
	opt := option.WithCredentialsFile("./bit-bot-188313-f3427d1a8526.json")
	config := &firebase.Config{
		DatabaseURL: "https://bit-bot-188313.firebaseio.com",
	}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
		return err
	}

	db, err := app.Database(ctx)
	if err != nil {
		log.Error(err)
	}

	path := "bitflyer/" + where
	ref := db.NewRef(path)
	if _, err := ref.Push(ctx, o); err != nil {
		log.Errorf("Failed adding alovelace: %v", err)
	}

	return nil
}
