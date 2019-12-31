package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestFetchINESHeader(t *testing.T) {
	tests := []struct {
		name        string
		openRom     func() ([]byte, error)
		want        *INESHeader
		makeWantErr func() error
	}{
		{
			name: "when rom is valid, return header",
			openRom: func() ([]byte, error) {
				romPath := "../test/roms/hello-world/hello-world.nes"
				f, err := os.Open(romPath)
				if err != nil {
					return nil, fmt.Errorf("failed to open rom\nromPath: %#v\nerr: %w", romPath, err)
				}
				defer f.Close()

				b, err := ioutil.ReadAll(f)
				if err != nil {
					return nil, fmt.Errorf("failed to open rom\nromPath: %#v\nerr: %w", romPath, err)
				}
				return b, nil
			},
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
				return fmt.Errorf("failed to fetch, rom is nil")
			},
		},
		{
			name:    "when rom is too short, return error",
			openRom: func() ([]byte, error) { return []byte{0x00, 0x01, 0x02, 0x03, 0x04}, nil },
			want:    nil,
			makeWantErr: func() error {
				return fmt.Errorf("failed to fetch, rom is too short")
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

			got, err := fetchINESHeader(rom)
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
