package view

import (
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
)

func TestNewSuccessResponse(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "correct",
			args: args{
				data: "test",
			},
			want: &Response{
				Data: "test",
			},
		},
		{
			name: "correct struct",
			args: args{
				data: struct {
					Name string `json:"name"`
				}{
					Name: "test",
				},
			},
			want: &Response{
				Data: struct {
					Name string `json:"name"`
				}{
					Name: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSuccessResponse(tt.args.data); !assert.Equal(t, got, tt.want) {
				t.Errorf("NewOKResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewErrorResponse(t *testing.T) {
	type args struct {
		errTest diterrors.StringError
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "correct",
			args: args{
				errTest: "some error text",
			},
			want: &Response{Error: &ErrorResponse{Message: "some error text"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrorResponse(tt.args.errTest); !assert.Equal(t, got, tt.want) {
				t.Errorf("NewErrorResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
