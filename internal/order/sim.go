package order

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/Rx-11/hft-arbitrage-engine/internal/arbitrage"
)

type SimExecutor struct {
	TradeSizeUSD float64
	SlippagePct  float64
	FeesPct      float64
	CumPnL       float64
	Trades       int
	Log          *logrus.Logger
}

func (s *SimExecutor) Execute(_ context.Context, opp arbitrage.Opportunity) {
	if s.TradeSizeUSD <= 0 {
		return
	}

	qty := s.TradeSizeUSD / opp.BuyPrice

	buyPx := opp.BuyPrice * (1.0 + s.SlippagePct/100.0)
	sellPx := opp.SellPrice * (1.0 - s.SlippagePct/100.0)

	buyNotional := buyPx * qty
	sellNotional := sellPx * qty
	fees := (buyNotional + sellNotional) * (s.FeesPct / 100.0)
	pnl := sellNotional - buyNotional - fees

	s.CumPnL += pnl
	s.Trades++

	if s.Log != nil {
		s.Log.WithFields(logrus.Fields{
			"symbol":        opp.Symbol,
			"buy_exchange":  opp.BuyExchange,
			"sell_exchange": opp.SellExchange,
			"buy_price":     buyPx,
			"sell_price":    sellPx,
			"spread_pct":    opp.SpreadPct,
			"trade_qty":     qty,
			"fees":          fees,
			"trade_pnl":     pnl,
			"cum_pnl":       s.CumPnL,
			"trades":        s.Trades,
		}).Info("Simulated arbitrage executed")
	}
}
