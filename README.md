# GoTorrent

![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=flat&logo=go)
![Protocol](https://img.shields.io/badge/Protocol-BitTorrent%20%2F%20BEP--15-orange)

A high-performance, lightweight BitTorrent client built from scratch in Go. Designed complying with **BEP-15** standards.

---
## Features

* **Custom Bencode Parser:** Built entirely without external libraries to handle raw torrent metadata structures (strings, integers, lists, and nested dictionaries).

* **HTTP & UDP Trackers Support:** Dynamic connection and failover across both HTTP and UDP trackers defined in the `announce` fields,  with custom retry mechanisms for UDP.

* **TCP Handshake & Wire Protocol:** Native TCP handshake and binary message engine (`Keep-Alive`, `Choke`, `Unchoke`, `Interested`, `Have`, `Bitfield`, `Request`, `Piece`, `Cancel`).

* **Disk Resume & Integrity Checking:** On-startup disk scanner that verifies existing local pieces against SHA-1 torrent hashes to resume downloads seamlessly to avoid re-downloading pieces that have already been successfully downloaded.

* **Cross-Platform Memory Pre-allocation:** Fast, non-sparse disk space reservation for both **Linux** and **Windows** prior to downloading.

* **Multi-Tier Fault Tolerance:**
    * Level 1: Automatic piece re-queuing on hash verification failures.

    * Level 2: Dynamic peer banning and disconnect triggers for continuously  malformed packets or corrupted payload delivery.

* **Test-Driven Design:** Developed using TDD practices for critical components like decoding, encoding, binary packing etc.

---


## Technical Decisions & Trade-Offs
### 1. Lock-Free `workChannel` vs. Centralized Manager

* **Decision:** Uses a shared buffered Go channel (`workChannel`) for distributing pieces to peer workers without coarse-grained mutex locks while a worker searches for its next piece. Workers maintain a local `visited` map to skip unsuitable pieces.

* **Trade-Off:** Provides zero-lock contention on torrents. However, on large torrents with thousands of pieces, rotating skipped pieces through the channel under high peer contention can lead to CPU overhead (*livelock potential*).

* **Considered Alternatives:**
  * *Global Mutex:* Locking the piece search phase globally, which slows throughput as all idle peers wait sequentially for work retrieval.
  * *Centralized Piece Coordinator:* Replacing channel rotation with a centralized manager tracking bitmasks of requested/downloaded pieces for lock-free bitfield lookups.

### 2. Sequential Downloading vs. Rarest-First
* **Decision:** Implements sequential piece ordering by default.

* **Trade-Off:** Good for media streaming/previewing (e.g. watching video content while downloading), but sub-optimal for overall swarm health compared to standard BitTorrent rarest-first piece distribution.

---

## Quick Start
### Prerequisites

* **Go 1.21** or higher installed on your system.

### Installation & Build

1. **Clone the repository:**

    ```bash
    # Clone the repository
    git clone https://github.com/Petroviiic/GoTorrent 
    ```

2. **Build & Run**

    ```bash
   # Direct execution
    go run ./... .\sample.torrent
   ``` 
   
> **Note:** The download destination is currently configured via the `DOWNLOAD_PATH` constant in `main.go` (defaults to `E:\GoBittorrentClient`). Make sure to adjust this path to a valid directory on your system before running.

## Future TO-DOs / Missing Features

* **Rarest-First Strategy:** Implement dynamic swarm piece frequency tracking to prioritize rare pieces.
* **Centralized Piece Coordinator:** Replace channel rotation with a centralized scheduler.
* **Seeding Support:** Upload verified pieces to requesting peers in the swarm.