package converter

import (
	"reflect"
	"testing"
)

func TestBytesToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "correct",
			args: args{
				b: []byte{
					116, 101, 115, 116, 32,
					115, 116, 114, 105, 110, 103, 32,
					116, 111, 32,
					99, 111, 110, 118, 101, 114, 116, 32,
					116, 111, 32,
					98, 121, 116, 101, 115, 32,
					230, 151, 165, 230, 156, 172, 228, 186, 186, 32,
					228, 184, 173, 229, 156, 139, 231, 154, 132,
				},
			},
			want: "test string to convert to bytes 日本人 中國的",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToString(tt.args.b); got != tt.want {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "correct",
			args: args{s: "test string to convert to bytes 日本人 中國的"},
			want: []byte{
				116, 101, 115, 116, 32,
				115, 116, 114, 105, 110, 103, 32,
				116, 111, 32,
				99, 111, 110, 118, 101, 114, 116, 32,
				116, 111, 32,
				98, 121, 116, 101, 115, 32,
				230, 151, 165, 230, 156, 172, 228, 186, 186, 32,
				228, 184, 173, 229, 156, 139, 231, 154, 132,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
