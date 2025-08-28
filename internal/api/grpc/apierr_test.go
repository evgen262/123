package grpc

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewApiError(t *testing.T) {
	type args struct {
		code codes.Code
		msg  string
		errs []error
	}
	testErr := []error{fmt.Errorf("testErr")}

	tests := []struct {
		name string
		args args
		want *apiError
	}{
		{
			name: "correct without errs",
			args: args{
				code: 502,
				msg:  "testMsg",
			},
			want: &apiError{
				Code:    502,
				Message: "testMsg",
			},
		},
		{
			name: "correct with errs",
			args: args{
				code: 502,
				msg:  "testMsg",
				errs: testErr,
			},
			want: &apiError{
				Code:    502,
				Message: "testMsg",
				Err:     &multierror.Error{Errors: testErr},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			err := NewApiError(tt.args.code, tt.args.msg, tt.args.errs...)
			if want.Err != nil && len(want.Err.Errors) > 0 {
				want.Err.ErrorFormat = nil
				err.Err.ErrorFormat = nil
			}
			assert.Equal(t, want, err)
		})
	}
}

func Test_apiError_Error(t *testing.T) {
	type fields struct {
		Code    codes.Code
		Message string
		Err     *multierror.Error
	}
	testErr := fmt.Errorf("testErr")

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "correct with undefined msg",
			fields: fields{
				Code: codes.Internal,
			},
			want: "error message not specified",
		},
		{
			name: "correct only msg",
			fields: fields{
				Code:    codes.Internal,
				Message: "testMsg",
			},
			want: "testMsg",
		},
		{
			name: "correct only err",
			fields: fields{
				Code: codes.Internal,
				Err:  &multierror.Error{Errors: []error{testErr}},
			},
			want: "testErr",
		},
		{
			name: "correct msg and err",
			fields: fields{
				Code:    codes.Internal,
				Message: "testMsg",
				Err:     &multierror.Error{Errors: []error{testErr}},
			},
			want: "testMsg: testErr",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var errString string
			if tt.fields.Err != nil {
				errString = NewApiError(tt.fields.Code, tt.fields.Message, tt.fields.Err).Error()
			} else {
				errString = NewApiError(tt.fields.Code, tt.fields.Message).Error()
			}
			assert.Equal(t, tt.want, errString)
		})
	}
}

func Test_apiError_GRPCStatus(t *testing.T) {
	type args struct {
		Code            codes.Code
		Message         string
		localizeMessage string
		Details         []proto.Message
	}
	testMsg := "testMsg"
	details := []proto.Message{&errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       "testField",
				Description: "testDesc",
			},
		},
	}}
	localizeMsg := &errdetails.LocalizedMessage{
		Locale:  "ru-RU",
		Message: "тестовое",
	}
	statusErrWithDetails, _ := status.New(codes.Internal, testMsg).WithDetails(details...)
	statusErrWithLocalize, _ := status.New(codes.Internal, testMsg).WithDetails(localizeMsg)

	tests := []struct {
		name string
		args args
		want *status.Status
	}{
		{
			name: "correct with empty message",
			args: args{
				Code:    codes.Internal,
				Message: "",
			},
			want: status.New(codes.Internal, "Что-то пошло не так. Попробуйте еще раз."),
		},
		{
			name: "correct without details",
			args: args{
				Code:    codes.Internal,
				Message: testMsg,
			},
			want: status.New(codes.Internal, testMsg),
		},
		{
			name: "correct with details",
			args: args{
				Code:    codes.Internal,
				Message: testMsg,
				Details: details,
			},
			want: statusErrWithDetails,
		},
		{
			name: "correct with localized",
			args: args{
				Code:            codes.Internal,
				Message:         testMsg,
				localizeMessage: localizeMsg.GetMessage(),
			},
			want: statusErrWithLocalize,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			err := NewApiError(tt.args.Code, tt.args.Message).
				WithDetails(tt.args.Details...).
				WithLocalizedMessage(tt.args.localizeMessage).
				GRPCStatus()
			assert.Equal(t, want, err)
		})
	}
}

func Test_apiError_WithDetails(t *testing.T) {
	type fields struct {
		Code    codes.Code
		Message string
	}
	type args struct {
		details []proto.Message
	}
	details := []proto.Message{&errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       "testField",
				Description: "testDesc",
			},
		},
	}}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *apiError
	}{
		{
			name: "correct",
			args: args{
				details: details,
			},
			want: func(a args, f fields) *apiError {
				return &apiError{
					Code:    f.Code,
					Message: f.Message,
					Details: a.details,
				}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				f := fields{
					Code:    codes.Internal,
					Message: "testMsg",
				}
				want := tt.want(tt.args, f)
				err := NewApiError(f.Code, f.Message).
					WithDetails(tt.args.details...)
				assert.Equal(t, want, err)
			})
		})
	}
}

func Test_apiError_WithLocalizedMessage(t *testing.T) {
	type args struct {
		msg string
	}
	type fields struct {
		Code    codes.Code
		Message string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) *apiError
	}{
		{
			name: "correct with localized",
			args: args{
				msg: "тестовое",
			},
			want: func(a args, f fields) *apiError {
				localizeMsg := &errdetails.LocalizedMessage{
					Locale:  "ru-RU",
					Message: a.msg,
				}
				return &apiError{
					Code:            f.Code,
					Message:         f.Message,
					localizeMessage: localizeMsg,
				}
			},
		},
		{
			name: "correct with empty localized msg",
			args: args{
				msg: "",
			},
			want: func(a args, f fields) *apiError {
				return &apiError{
					Code:    f.Code,
					Message: f.Message,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				f := fields{
					Code:    codes.Internal,
					Message: "testMsg",
				}
				want := tt.want(tt.args, f)
				err := NewApiError(f.Code, f.Message).
					WithLocalizedMessage(tt.args.msg)
				assert.Equal(t, want, err)
			})
		})
	}
}
