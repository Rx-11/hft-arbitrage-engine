# HFT Arbitrage Engine

A lightweight, modular arbitrage engine written in Go.  
It connects to **Finnhub’s WebSocket API** for real-time market data,  
detects **cross-exchange arbitrage opportunities**, and simulates trades with PnL tracking.

---

## Features

- Real-time market data streaming via Finnhub WebSocket
- Multi-exchange crypto pair support (Binance, Coinbase, Kraken, Bitstamp, Gemini, etc.)
- Cross-exchange arbitrage detection with configurable spread threshold and staleness window
- Simulated trade execution with slippage and fees
- Structured JSON logging using `logrus`
- Modular, production-friendly project layout

---

## Project Structure

```
internal/
  websocket/   # WebSocket connector for Finnhub
  handler/     # Parses and normalizes incoming market data
  arbitrage/   # Arbitrage detection engine
  order/       # Simulated trade executor
```

---

## Requirements

- [Go 1.20+](https://go.dev/dl/)
- [Finnhub API key](https://finnhub.io/register)

---

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/hft-arbitrage-engine.git
   cd hft-arbitrage-engine
   ```

2. Create a .env file in the project root:
```bash
  API_KEY=your_finnhub_api_key
  LOG_LEVEL=info
```

3. Or export your Finnhub API key:
   ```bash
   export FINNHUB_API_KEY=your_api_key_here
   export LOG_LEVEL=debug
   ```

4. Run the engine:
   ```bash
   go run main.go
   ```

---

## Configuration

Edit the symbol list in `main.go`:
```go
symbols := []string{
  "BINANCE:BTCUSDT",
  "COINBASE:BTC-USD",
  "KRAKEN:XBTUSD",
  "BITSTAMP:BTCUSD",
}
```

Adjust trading parameters:
```go
thresholdPct := 0.10             // Minimum spread % to trigger arbitrage
staleAfter   := 2 * time.Second  // Ignore stale quotes older than this
tradeSizeUSD := 1000.0           // Simulated trade size
slippagePct  := 0.01             // Slippage per leg
feesPct      := 0.02             // Fees per leg
```

---

## Example Output

```json
{"exchange":"BINANCE","symbol":"BTC/USD","price":27234.5,"level":"debug","msg":"processed tick"}
{"symbol":"BTC/USD","buy_exchange":"COINBASE","sell_exchange":"BINANCE","spread_pct":0.15,"level":"info","msg":"Arbitrage opportunity detected"}
{"symbol":"BTC/USD","buy_ex":"COINBASE","sell_ex":"BINANCE","trade_pnl_usd":1.20,"cum_pnl_usd":1.20,"level":"info","msg":"Simulated arbitrage executed"}
```

---

## Disclaimer

This software is for educational and research purposes only.  
It does not provide financial advice. Use responsibly.
