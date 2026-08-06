# Lagless Relay

Go UDP relay for latency-sensitive traffic — the networking core of an edge latency / predictive-routing experiment.

## What it does
- Listens for sequenced UDP packets from clients
- Forwards payloads to a local game/application server
- Sends client acknowledgements on a dedicated port
- Logs forward/ack throughput for operational debugging

This repository currently publishes the **Go relay service**. Related edge-ML work (ONNX latency forecasting, TUN mesh) is documented on LinkedIn/resume; contributions that land here will be reflected in this README.

## Quick start

```bash
go run .

# Optional env
# LAGLESS_LISTEN=0.0.0.0:7001   # client-facing listen address
# LAGLESS_GAME=127.0.0.1:7002   # upstream game/app UDP address
```

## Stack
- Go 1.22+
- UDP (no external dependencies)

## Status
v0.1 — functional relay prototype suitable for local experimentation and load testing.

## License
Add a license if you intend this to be reused publicly (MIT recommended).
