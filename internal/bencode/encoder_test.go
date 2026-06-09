package bencode_test

import (
	"testing"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
	"github.com/google/go-cmp/cmp"
)

func TestEncoders(t *testing.T) {
	tests := []struct {
		name        string
		data        any
		wantErr     bool
		expectedRes []byte
	}{
		{
			name:        "longer word encoding",
			data:        "computerscience",
			wantErr:     false,
			expectedRes: []byte("15:computerscience"),
		},

		{
			name:        "integer encoding; positive number",
			data:        0,
			wantErr:     false,
			expectedRes: []byte("i0e"),
		},
		{
			name:        "integer encoding; negative number",
			data:        -42,
			wantErr:     false,
			expectedRes: []byte("i-42e"),
		},

		{
			name:        "list encoding 1",
			data:        []any{"bencode"},
			wantErr:     false,
			expectedRes: []byte("l7:bencodee"),
		},
		{
			name:        "list encoding 2",
			data:        []any{"bencode", -20},
			wantErr:     false,
			expectedRes: []byte("l7:bencodei-20ee"),
		},

		{
			name:        "dictionary encoding",
			data:        map[any]any{"meaning": 42, "wiki": "bencode"},
			wantErr:     false,
			expectedRes: []byte("d7:meaningi42e4:wiki7:bencodee"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := bencode.Encode(tt.data)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("encoder failed: %v", err)
				}
				return
			}

			if diff := cmp.Diff(tt.expectedRes, res); diff != "" {
				t.Errorf("encoder failed. Wrong result mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
