package tracker_test

import (
	"net"
	"testing"

	"github.com/Petroviiic/GoTorrent/internal/tracker"
	"github.com/google/go-cmp/cmp"
)

func TestDecodePeerList(t *testing.T) {
	tests := []struct {
		name  string
		peers []byte
		want  []*tracker.Peer
	}{
		{
			name:  "Empty peers byte list - no peers",
			peers: []byte{},
			want:  []*tracker.Peer{},
		},
		{
			name:  "One valid peer (192.168.1.1:6881)",
			peers: []byte{192, 168, 1, 1, 0x1A, 0xE1},
			want: []*tracker.Peer{
				{
					IP:   net.ParseIP("192.168.1.1"),
					Port: 6881,
				},
			},
		},
		{
			name: "Multiple valid peers",
			peers: []byte{
				127, 0, 0, 1, 0x1F, 0x90,
				8, 8, 8, 8, 0x00, 0x35,
			},
			want: []*tracker.Peer{
				{
					IP:   net.ParseIP("127.0.0.1"),
					Port: 8080,
				},
				{
					IP:   net.ParseIP("8.8.8.8"),
					Port: 53,
				},
			},
		},
		{
			name:  "Invalid peer list size 1",
			peers: []byte{192, 168, 1, 1},
			want:  []*tracker.Peer{},
		},
		{
			name: "Invalid peer list size 2",
			peers: []byte{
				192, 168, 1, 1, 0x1F, 0x90,
				8, 8, 8,
			},
			want: []*tracker.Peer{
				{
					IP:   net.ParseIP("192.168.1.1"),
					Port: 8080,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tracker.DecodePeerList(tt.peers)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Peer list decoder failed. Wrong result mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
