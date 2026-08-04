# GoTorrent

A high-performance, lightweight BitTorrent client built from scratch in Go. All features implemented in this project follow BEP-15 protocol rules.

## Features

*    **Custom Bencode Parser:** Built without external libraries to handle all raw torrent file structures, such as strings, integers, lists, and nested dictionaries.

*   **HTTP & UDP Trackers Support:** Connection to both HTTP and UDP trackers gathered from announce field of torrent file, with dynamic failover and retry mechanism for UDP trackers.

*    **Peer Wire Protocol:** Native TCP handshake and binary message engine (Keep-Alive, Choke, Unchoke, Interested, Have, Bitfield, Request, Piece, Cancel).

*   **Disk Resume & Integrity Checking:** On-startup disk scanner that verifies existing local pieces against SHA-1 torrent hashes to resume downloads seamlessly and avoid re-downloading pieces that have already been successfully downloaded.

*    **Cross-Platform Memory Pre-allocation:** Fast, non-sparse disk space reservation for both Linux and Windows prior to downloading.


