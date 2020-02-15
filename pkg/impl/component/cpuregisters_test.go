package component_test

import "testing"

import "nes-go/pkg/impl/component"

func TestCPUStatusRegisterUpdateN(t *testing.T) {
	tests := []struct {
		name string
		ans  byte
		want bool
	}{
		{
			name: "When ans is 0x00, Negative is false",
			ans:  0x00,
			want: false,
		},
		{
			name: "When ans is 0x7F, Negative is false",
			ans:  0x7F,
			want: false,
		},
		{
			name: "When ans is 0x80, Negative is true",
			ans:  0x80,
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := component.NewCPUStatusRegister()
			r.UpdateN(test.ans)
			got := r.Negative
			if got != test.want {
				t.Errorf("wrong parameter\nwant:%v\ngot :%v", test.want, got)
			}
		})
	}
}

func TestCPUStatusRegisterUpdateV(t *testing.T) {
	tests := []struct {
		name   string
		before byte
		after  byte
		want   bool
	}{
		{
			name:   "When value changed positive to positive, Overflow is false",
			before: 0x7E,
			after:  0x7F,
			want:   false,
		},
		{
			name:   "When value changed positive to negative, Overflow is true",
			before: 0x7F,
			after:  0xFF,
			want:   true,
		},
		{
			name:   "When value changed negative to positive, Overflow is true",
			before: 0x80,
			after:  0x7F,
			want:   true,
		},
		{
			name:   "When value changed negative to negative, Overflow is false",
			before: 0x80,
			after:  0xFF,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := component.NewCPUStatusRegister()
			r.UpdateV(test.before, test.after)
			got := r.Overflow
			if got != test.want {
				t.Errorf("wrong parameter\nvalue:%#2X=>%#2X\noverflow want:%v\noverflow got :%v", test.before, test.after, test.want, got)
			}
		})
	}
}

func TestCPUStatusRegisterUpdateC(t *testing.T) {
	tests := []struct {
		name string
		ans  uint16
		want bool
	}{
		{
			name: "When ans <= 0x00FF, Carry is false",
			ans:  0x00FF,
			want: false,
		},
		{
			name: "When ans >= 0x0100, Carry is true",
			ans:  0x0100,
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := component.NewCPUStatusRegister()
			r.UpdateC(test.ans)
			got := r.Carry
			if got != test.want {
				t.Errorf("wrong parameter\nans:%#4X\ncarry want:%v\ncarry got :%v", test.ans, test.want, got)
			}
		})
	}
}
