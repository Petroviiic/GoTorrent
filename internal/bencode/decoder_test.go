package bencode_test

import (
	"testing"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
	"github.com/google/go-cmp/cmp"
)

func TestDecoders(t *testing.T) {

	tests := []struct {
		name          string
		buffer        []byte
		wantErr       bool
		expectedRes   any
		expectedIndex int
	}{
		{
			name:          "short word decoding",
			buffer:        []byte("5:world"),
			wantErr:       false,
			expectedRes:   "world",
			expectedIndex: len([]byte("5:world")),
		},
		{
			name:          "longer word decoding",
			buffer:        []byte("15:computerscience"),
			wantErr:       false,
			expectedRes:   "computerscience",
			expectedIndex: len([]byte("15:computerscience")),
		},

		{
			name:          "integer decoding; positive number",
			buffer:        []byte("i0e"),
			wantErr:       false,
			expectedRes:   0,
			expectedIndex: len([]byte("i0e")),
		},
		{
			name:          "integer decoding; negative number",
			buffer:        []byte("i-42e"),
			wantErr:       false,
			expectedRes:   -42,
			expectedIndex: len([]byte("i-42e")),
		},
		{
			name:          "integer decoding; extra digits",
			buffer:        []byte("i42e4"),
			wantErr:       false,
			expectedRes:   42,
			expectedIndex: len([]byte("i42e")),
		},

		{
			name:          "list decoding 1",
			buffer:        []byte("l7:bencodee"),
			wantErr:       false,
			expectedRes:   []any{"bencode"},
			expectedIndex: len([]byte("l7:bencodee")),
		},
		{
			name:          "list decoding 2",
			buffer:        []byte("l7:bencodei-20ee"),
			wantErr:       false,
			expectedRes:   []any{"bencode", -20},
			expectedIndex: len([]byte("l7:bencodei-20ee")),
		},

		{
			name:          "dictionary decoding",
			buffer:        []byte("d7:meaningi42e4:wiki7:bencodee"),
			wantErr:       false,
			expectedRes:   map[any]any{"meaning": 42, "wiki": "bencode"},
			expectedIndex: len([]byte("d7:meaningi42e4:wiki7:bencodee")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := bencode.NewDecoder(tt.buffer)

			res, index, err := decoder.Decoders[tt.buffer[0]](0)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("Decoder failed: %v", err)
				}
				return
			}

			if index != tt.expectedIndex {
				t.Errorf("Decoder failed. Wrong index: expected %v , got %v", tt.expectedIndex, index)
				return
			}

			if diff := cmp.Diff(tt.expectedRes, res); diff != "" {
				t.Errorf("Decoder failed. Wrong result mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
