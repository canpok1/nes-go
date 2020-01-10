package domain

import "testing"

func TestNameTablePointToAttributeTableIndex(t *testing.T) {
	tests := []struct {
		name string
		x    uint8
		y    uint8
		want uint16
	}{
		{
			name: "when point is 1st col and 1st row, return 0",
			x:    0,
			y:    0,
			want: 0,
		},
		{
			name: "when point is 2nd col and 1st row, return 0",
			x:    1,
			y:    0,
			want: 0,
		},
		{
			name: "when point is 3rd col and 1st row, return 0",
			x:    2,
			y:    0,
			want: 0,
		},
		{
			name: "when point is 5th col and 1st row, return 1",
			x:    4,
			y:    0,
			want: 1,
		},
		{
			name: "when point is last col and 1st row, return 7",
			x:    31,
			y:    0,
			want: 7,
		},
		{
			name: "when point is 1st col and 2nd row, return 0",
			x:    0,
			y:    1,
			want: 0,
		},
		{
			name: "when point is 1st col and 3rd row, return 0",
			x:    0,
			y:    2,
			want: 0,
		},
		{
			name: "when point is 1st col and 5th row, return 8",
			x:    0,
			y:    4,
			want: 8,
		},
		{
			name: "when point is 1st col and last row, return 56",
			x:    0,
			y:    29,
			want: 56,
		},
		{
			name: "when point is last col and last row, return 63",
			x:    31,
			y:    29,
			want: 63,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NameTablePoint{X: tt.x, Y: tt.y}
			got := p.ToAttributeTableIndex()
			if got != tt.want {
				t.Errorf("wrong output\nwant: %#v\ngot: %#v", tt.want, got)
			}
		})
	}
}
