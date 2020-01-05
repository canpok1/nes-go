package model

import "testing"

func TestToAttributeTableIndex(t *testing.T) {
	tests := []struct {
		name     string
		x        uint8
		y        uint8
		want     int
		hasError bool
	}{
		{
			name:     "when palette position is first col and first row, return 0",
			x:        0,
			y:        0,
			want:     0,
			hasError: false,
		},
		{
			name:     "when palette position is second col and first row, return 0",
			x:        16,
			y:        0,
			want:     0,
			hasError: false,
		},
		{
			name:     "when palette position is third col and first row, return 1",
			x:        32,
			y:        0,
			want:     1,
			hasError: false,
		},
		{
			name:     "when palette position is last col and first row, return 7",
			x:        240,
			y:        0,
			want:     7,
			hasError: false,
		},
		{
			name:     "when palette position is first col and second row, return 0",
			x:        0,
			y:        16,
			want:     0,
			hasError: false,
		},
		{
			name:     "when palette position is first col and third row, return 8",
			x:        0,
			y:        32,
			want:     8,
			hasError: false,
		},
		{
			name:     "when palette position is first col and last row, return 48",
			x:        0,
			y:        208,
			want:     48,
			hasError: false,
		},
		{
			name:     "when palette position is last col and last row, return 55",
			x:        224,
			y:        208,
			want:     55,
			hasError: false,
		},
		{
			name:     "when y is out of range, return error",
			x:        0,
			y:        240,
			want:     0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toAttributeTableIndex(MonitorX(tt.x), MonitorY(tt.y))
			if tt.hasError != (err != nil) {
				t.Errorf("wrong error\nhasError: %v\nerror: %v", tt.hasError, err)
				return
			}
			if got != tt.want {
				t.Errorf("wrong output\nwant: %#v\ngot: %#v", tt.want, got)
			}
		})
	}
}
