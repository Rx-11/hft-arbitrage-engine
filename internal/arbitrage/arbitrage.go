package arbitrage

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Opportunity struct {
	Symbol       string
	BuyExchange  string
	SellExchange string
	BuyPrice     float64
	SellPrice    float64
	SpreadPct    float64
	Timestamp    time.Time
}

type Quote struct {
	Price float64
	Time  time.Time
}

type Detector struct {
	mu         sync.RWMutex
	quotes     map[string]map[string]Quote // symbol -> exchange -> quote
	threshold  float64
	staleAfter time.Duration
	log        *logrus.Logger
}

func NewDetector(thresholdPct float64, staleAfter time.Duration, log *logrus.Logger) *Detector {
	return &Detector{
		quotes:     make(map[string]map[string]Quote),
		threshold:  thresholdPct,
		staleAfter: staleAfter,
		log:        log,
	}
}

func (d *Detector) UpdateQuote(symbol, exchange string, price float64, ts time.Time) *Opportunity {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.quotes[symbol]; !ok {
		d.quotes[symbol] = make(map[string]Quote)
	}
	d.quotes[symbol][exchange] = Quote{Price: price, Time: ts}

	return d.checkForArbitrage(symbol)
}

func (d *Detector) checkForArbitrage(symbol string) *Opportunity {
	now := time.Now()
	quotes := d.quotes[symbol]

	var (
		bestBuyEx, bestSellEx string
		bestBuy, bestSell     float64
		found                 bool
	)

	for ex, q := range quotes {
		if now.Sub(q.Time) > d.staleAfter {
			continue
		}
		if !found {
			bestBuy, bestSell = q.Price, q.Price
			bestBuyEx, bestSellEx = ex, ex
			found = true
			continue
		}
		if q.Price < bestBuy {
			bestBuy = q.Price
			bestBuyEx = ex
		}
		if q.Price > bestSell {
			bestSell = q.Price
			bestSellEx = ex
		}
	}

	if !found || bestBuyEx == bestSellEx {
		return nil
	}

	spreadPct := (bestSell - bestBuy) / bestBuy * 100
	if spreadPct >= d.threshold {
		opp := &Opportunity{
			Symbol:       symbol,
			BuyExchange:  bestBuyEx,
			SellExchange: bestSellEx,
			BuyPrice:     bestBuy,
			SellPrice:    bestSell,
			SpreadPct:    spreadPct,
			Timestamp:    now,
		}
		d.log.WithFields(logrus.Fields{
			"symbol":        symbol,
			"buy_exchange":  bestBuyEx,
			"sell_exchange": bestSellEx,
			"buy_price":     bestBuy,
			"sell_price":    bestSell,
			"spread_pct":    spreadPct,
		}).Info("Arbitrage opportunity detected")
		return opp
	}

	return nil
}
