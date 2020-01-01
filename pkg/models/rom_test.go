package model

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
