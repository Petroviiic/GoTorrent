# GoTorrent

A high-performance, lightweight BitTorrent client built from scratch in Go. All features implemented in this project follow BEP-15 protocol rules.

## Features

*    **Custom Bencode Parser:** Built without external libraries to handle all raw torrent file structures, such as strings, integers, lists, and nested dictionaries.

*   **HTTP & UDP Trackers Support:** Connection to both HTTP and UDP trackers gathered from announce field of torrent file, with dynamic failover and retry mechanism for UDP trackers.

*    **TCP Handshake & Messaging:** Native TCP handshake and binary message engine (Keep-Alive, Choke, Unchoke, Interested, Have, Bitfield, Request, Piece, Cancel).

*   **Disk Resume & Integrity Checking:** On-startup disk scanner that verifies existing local pieces against SHA-1 torrent hashes to resume downloads seamlessly and avoid re-downloading pieces that have already been successfully downloaded.

*    **Cross-Platform Memory Pre-allocation:** Fast, non-sparse disk space reservation for both Linux and Windows prior to downloading.

* **Faulty Pieces Tolerance:**

    * Level 1: Automatic piece re-queuing on hash verification failures.

    * Level 2: Dynamic peer banning and disconnect triggers for continiously  malformed packets or corrupted payload delivery.

* **Test-Driven Design:** Developed using TDD practices for critical path components like parsing, binary packing, and hashing.


<br>
<br>

## Technical Decisions & Trade-Offs
1. Lock-Free workChannel vs. Centralized Manager

    Decision: Current implementation uses a shared buffered Go channel for distributing pieces to peer workers without mutex locking while one worker searches for its next piece. Peer workers use a local visited map to skip unsuitable pieces.

    Trade-Off: Provides zero-lock contention on torrents. However, on large torrents with thousands of pieces, rotating skipped pieces through the channel can potentially lead to livelock under high peer contention.

    Other options: 
    * Adding a global mutex for while the peer is seaching for a new piece, which could slow down the whole process, because all peers that need a new piece would need to wait for the first one to retrieve the next piece of work.
    * Replacing channel rotation with a centralized manager that would keep track of (bitmasks for) all requested and downloaded pieces, so that workers can request new piece based on their bitfield. 

2. Sequential Downloading vs. Rarest-First

    Decision: Implemented sequential piece ordering by default.

    Trade-Off: Ideal for media streaming/previewing (e.g., watching video while downloading), but sub-optimal for overall swarm health compared to standard BitTorrent rarest-first piece distribution.
