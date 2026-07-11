## GoTorrent

A lightweight, concurrent BitTorrent client built from scratch in Go. 

## Features
* **Custom Protocol Implementation:** Built-in network message serializer/deserializer and handshake verification.
* **Concurrent Downloader:** Uses a thread-safe worker pool model to download multiple file pieces in parallel.
* **Asynchronous Pipelining:** Efficiently requests 16KB blocks concurrently from unchoked peers.

## Tech Stack
* **Language:** Go (Golang)
* **Packages:** `net` (TCP communication), `crypto/sha1` (data integrity verification).