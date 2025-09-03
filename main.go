package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Rx-11/hft-arbitrage-engine/config"
	"github.com/Rx-11/hft-arbitrage-engine/internal/arbitrage"
	"github.com/Rx-11/hft-arbitrage-engine/internal/handler"
	"github.com/Rx-11/hft-arbitrage-engine/internal/order"
	"github.com/Rx-11/hft-arbitrage-engine/internal/websocket"
)

func main() {
	config.LoadConfig()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	if lvl, err := logrus.ParseLevel(config.Cfg.LogLevel); err == nil {
		log.SetLevel(lvl)
	}
	handler.SetLogger(log)

	symbols := []string{
		"BINANCE:BTCUSDT",
		"COINBASE:BTC-USD",
		"KRAKEN:XBTUSD",
		"BITSTAMP:BTCUSD",
		"GEMINI:BTCUSD",
		"OKX:BTC-USD",
	}

	tickCh := make(chan handler.Tick, 1000)
	handler.SetTickSink(tickCh)

	det := arbitrage.NewDetector(0.01, 2*time.Second, log)

	exec := &order.SimExecutor{
		TradeSizeUSD: 1000.0,
		SlippagePct:  0.001,
		FeesPct:      0.002,
		Log:          log,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for tk := range tickCh {
			if opp := det.UpdateQuote(tk.Symbol, tk.Source, tk.Price, tk.Time); opp != nil {
				exec.Execute(ctx, *opp)
			}
		}
	}()

	client := websocket.NewWebSocketClient(config.Cfg.ApiKey, symbols)
	if err := client.Connect(); err != nil {
		log.WithError(err).Fatal("WebSocket connect failed")
	}
	defer client.Close()

	client.Subscribe()
	go client.Listen()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
	close(tickCh)
	log.Info("Shutdown complete")
}
