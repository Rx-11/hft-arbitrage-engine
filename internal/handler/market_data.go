package handler

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func SetLogger(l *logrus.Logger) {
	if l != nil {
		log = l
	}
}

type Tick struct {
	Symbol string
	Price  float64
	Time   time.Time
	Source string
}

var (
	tickSink   chan<- Tick
	tickSinkMu sync.RWMutex
)

func SetTickSink(ch chan<- Tick) {
	tickSinkMu.Lock()
	tickSink = ch
	tickSinkMu.Unlock()
}

func sendTick(t Tick) {
	tickSinkMu.RLock()
	defer tickSinkMu.RUnlock()
	if tickSink == nil {
		return
	}
	select {
	case tickSink <- t:
	default:
		log.WithField("symbol", t.Symbol).Warn("tick dropped: sink full")
	}
}

type MarketData struct {
	Type string `json:"type"`
	Data []struct {
		S string  `json:"s"`
		P float64 `json:"p"`
		T int64   `json:"t"`
	} `json:"data"`
}

func ProcessMarketData(message []byte) {
	b := bytes.TrimSpace(message)
	if len(b) == 0 || b[0] != '{' {
		log.WithField("payload", string(b)).Debug("non-JSON frame; skipping")
		return
	}

	var md MarketData
	if err := json.Unmarshal(b, &md); err != nil {
		log.WithError(err).WithField("raw", string(b)).Warn("failed to parse JSON market data")
		return
	}

	if md.Type != "trade" {
		return
	}
	for _, d := range md.Data {
		ex, norm := splitExchangeSymbol(d.S)
		ts := time.UnixMilli(d.T)
		if d.T == 0 {
			ts = time.Now()
		}

		sendTick(Tick{
			Symbol: norm,
			Price:  d.P,
			Time:   ts,
			Source: ex,
		})

		log.WithFields(logrus.Fields{
			"exchange": ex,
			"symbol":   norm,
			"raw_sym":  d.S,
			"price":    d.P,
			"ts_ms":    d.T,
		}).Debug("processed tick")
	}
}

func splitExchangeSymbol(s string) (exchange string, normalized string) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "UNKNOWN", s
	}
	ex, raw := parts[0], parts[1]

	raw = strings.ReplaceAll(raw, "-", "")
	raw = strings.ReplaceAll(raw, "XBT", "BTC")

	for _, q := range []string{"USDT", "USD", "USDC", "BUSD"} {
		if strings.HasSuffix(raw, q) {
			base := strings.TrimSuffix(raw, q)
			// return ex, base + "/" + q
			return ex, base + "/USD"
		}
	}

	return ex, raw
}
