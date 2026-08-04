# GoTorrent

A high-performance, lightweight BitTorrent client built from scratch in Go. 

## Features

*    **Custom Bencode Parser:** Built without external libraries to handle raw torrent metadata, strings, integers, lists, and nested dictionaries.

*   **Multi-Tracker Architecture:** Supports both HTTP and UDP trackering (BEP 15) with dynamic failover between announce and announce-list endpoints.

*    **Full Peer Wire Protocol:** Native TCP handshake and binary message engine (Keep-Alive, Choke, Unchoke, Interested, Have, Bitfield, Request, Piece, Cancel).

*   **Disk Resume & Integrity Checking:** On-startup disk scanner that verifies existing local pieces against SHA-1 torrent hashes to resume downloads seamlessly.

*    **Cross-Platform Pre-allocation:** Fast, non-sparse disk space reservation for both Linux and Windows prior to downloading.

*    **Sequential Streaming Strategy:** In-order piece distribution via worker channels, enabling early file preview/streaming for media formats.

