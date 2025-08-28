package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_detailsError_Error(t *testing.T) {
	type args struct {
		detailsErr *DetailsError
	}

	testUUID := uuid.New()
	tests := []struct {
		name string
		args args
		want func(a args) string
	}{
		{
			name: "nil",
			args: args{
				detailsErr: nil,
			},
			want: func(a args) string {
				return ""
			},
		},
		{
			name: "correct",
			args: args{
				detailsErr: func() *DetailsError {
					e := DetailsError{message: testUUID.String()}
					return &e
				}(),
			},
			want: func(a args) string {
				return testUUID.String()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want(tt.args)
			e := tt.args.detailsErr
			got := e.Error()

			assert.Equal(t, want, got)
		})
	}
}

func Test_detailsError_GetField(t *testing.T) {
	type args struct {
		detailsErr *DetailsError
	}

	tests := []struct {
		name string
		args args
		want func(a args) string
	}{
		{
			name: "nil",
			args: args{
				detailsErr: nil,
			},
			want: func(a args) string {
				return ""
			},
		},
		{
			name: "correct",
			args: args{
				detailsErr: func() *DetailsError {
					e := DetailsError{field: "testField"}
					return &e
				}(),
			},
			want: func(a args) string {
				return "testField"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want(tt.args)
			e := tt.args.detailsErr
			got := e.GetField()

			assert.Equal(t, want, got)
		})
	}
}

func Test_detailsError_GetMessage(t *testing.T) {
	type args struct {
		detailsErr *DetailsError
	}

	testUUID := uuid.New()
	tests := []struct {
		name string
		args args
		want func(a args) string
	}{
		{
			name: "nil",
			args: args{
				detailsErr: nil,
			},
			want: func(a args) string {
				return ""
			},
		},
		{
			name: "correct",
			args: args{
				detailsErr: func() *DetailsError {
					e := DetailsError{message: testUUID.String()}
					return &e
				}(),
			},
			want: func(a args) string {
				return testUUID.String()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want(tt.args)
			e := tt.args.detailsErr
			got := e.GetMessage()

			assert.Equal(t, want, got)
		})
	}
}

func Test_detailsError_GetReauthRequired(t *testing.T) {
	type args struct {
		detailsErr *DetailsError
	}

	tests := []struct {
		name string
		args args
		want func(a args) bool
	}{
		{
			name: "nil",
			args: args{
				detailsErr: nil,
			},
			want: func(a args) bool {
				return false
			},
		},
		{
			name: "correct",
			args: args{
				detailsErr: func() *DetailsError {
					e := DetailsError{reauthRequired: true}
					return &e
				}(),
			},
			want: func(a args) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want(tt.args)
			e := tt.args.detailsErr
			got := e.GetReauthRequired()

			assert.Equal(t, want, got)
		})
	}
}

func Test_detailsError_NewDetailsError(t *testing.T) {
	type args struct {
		field, message string
		reauth         bool
	}

	tests := []struct {
		name string
		args args
		want *DetailsError
	}{
		{
			name: "correct",
			args: args{
				field:   "testField",
				message: "testMessage",
				reauth:  true,
			},
			want: &DetailsError{
				field:          "testField",
				message:        "testMessage",
				reauthRequired: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDetailsError(tt.args.field, tt.args.message, tt.args.reauth)

			assert.Equal(t, tt.want, got)
		})
	}
}
