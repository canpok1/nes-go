package domain

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParseINESHeader(t *testing.T) {
	tests := []struct {
		name        string
		openRom     func() ([]byte, error)
		want        *INESHeader
		makeWantErr func() error
	}{
		{
			name:    "when rom is valid, return header",
			openRom: func() ([]byte, error) { return openROM("../../test/roms/hello-world/hello-world.nes") },
			want: &INESHeader{
				PRGROMSize: 0x02,
				CHRROMSize: 0x01,
			},
			makeWantErr: func() error { return nil },
		},
		{
			name:    "when rom is nil, return error",
			openRom: func() ([]byte, error) { return nil, nil },
			want:    nil,
			makeWantErr: func() error {
				return fmt.Errorf("failed to parse, rom is nil")
			},
		},
		{
			name:    "when rom is too short, return error",
			openRom: func() ([]byte, error) { return []byte{0x00, 0x01, 0x02, 0x03, 0x04}, nil },
			want:    nil,
			makeWantErr: func() error {
				return fmt.Errorf("failed to parse, rom is too short")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rom, err := tt.openRom()
			if err != nil {
				t.Errorf("failed to open rom\nerr: %#v", err)
				return
			}

			wantErr := tt.makeWantErr()

			got, err := parseINESHeader(rom)
			if !reflect.DeepEqual(err, wantErr) {
				t.Errorf("wrong error\ngot: %#v\nwant: %#v", err, wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrong output\ngot: %#v\nwant: %#v", got, tt.want)
				return
			}
		})
	}
}

func TestFetchROM(t *testing.T) {
	tests := []struct {
		name     string
		romPath  string
		hasError bool
	}{
		{
			name:     "when rom is valid, not return error",
			romPath:  "../../test/roms/hello-world/hello-world.nes",
			hasError: false,
		},
	}

	for _, tt := range tests {
		if _, err := FetchROM(tt.romPath); (err != nil) != tt.hasError {
			t.Errorf("wrong error\ngot: %#v\nhasError: %#v", err, tt.hasError)
			return
		}
	}
}

func TestTilePatternToColorMap(t *testing.T) {
	tests := []struct {
		name string
		tile TilePattern
		want [][]byte
	}{
		{
			name: "when tile pattern is all 0, return all 0 colorMap",
			tile: TilePattern([]byte{
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000

				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
				0x00, // 0b 0000 0000
			}),
			want: [][]byte{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "when tile pattern is all 0xFF, return all 3 colorMap",
			tile: TilePattern([]byte{
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111

				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
			}),
			want: [][]byte{
				{3, 3, 3, 3, 3, 3, 3, 3},
				{3, 3, 3, 3, 3, 3, 3, 3},
				{3, 3, 3, 3, 3, 3, 3, 3},
				{3, 3, 3, 3, 3, 3, 3, 3},
				{3, 3, 3, 3, 3, 3, 3, 3},
				{3, 3, 3, 3, 3, 3, 3, 3},
				{3, 3, 3, 3, 3, 3, 3, 3},
				{3, 3, 3, 3, 3, 3, 3, 3},
			},
		},
		{
			// このページのパターンでテスト
			// https://qiita.com/bokuweb/items/1575337bef44ae82f4d3#%E3%82%AD%E3%83%A3%E3%83%A9%E3%82%AF%E3%82%BF%E3%83%BCrom
			name: "when tile pattern is hert, return hert pattern colorMap",
			tile: TilePattern([]byte{
				0x66, // 0b 0110 0110
				0x7F, // 0b 0111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0x7E, // 0b 0111 1110
				0x3C, // 0b 0011 1100
				0x18, // 0b 0001 1000

				0x66, // 0b 0110 0110
				0x5F, // 0b 0101 1111
				0xBF, // 0b 1011 1111
				0xBF, // 0b 1011 1111
				0xFF, // 0b 1111 1111
				0x7E, // 0b 0111 1110
				0x3C, // 0b 0011 1100
				0x18, // 0b 0001 1000
			}),
			want: [][]byte{
				{0, 3, 3, 0, 0, 3, 3, 0},
				{0, 3, 1, 3, 3, 3, 3, 3},
				{3, 1, 3, 3, 3, 3, 3, 3},
				{3, 1, 3, 3, 3, 3, 3, 3},
				{3, 3, 3, 3, 3, 3, 3, 3},
				{0, 3, 3, 3, 3, 3, 3, 0},
				{0, 0, 3, 3, 3, 3, 0, 0},
				{0, 0, 0, 3, 3, 0, 0, 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tile.toColorMap()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrong output\ntile: %#v\ngot: %#v\nwant: %#v", tt.tile, got, tt.want)
			}
		})
	}
}

func openROM(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to open rom\nromPath: %#v\nerr: %w", p, err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to open rom\nromPath: %#v\nerr: %w", p, err)
	}
	return b, nil
}
