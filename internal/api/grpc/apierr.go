package grpc

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type apiError struct {
	Code            codes.Code
	Message         string
	Err             *multierror.Error
	localizeMessage *errdetails.LocalizedMessage
	Details         []proto.Message
}

func NewApiError(code codes.Code, msg string, errs ...error) *apiError {
	var err *multierror.Error
	if len(errs) > 0 {
		err = &multierror.Error{
			ErrorFormat: func(errs []error) string {
				var stringErrs []string
				for _, err := range errs {
					stringErrs = append(stringErrs, err.Error())
				}
				return strings.Join(stringErrs, "; ")
			},
		}
		err = multierror.Append(err, errs...)
	}

	return &apiError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func (e *apiError) Error() string {
	if e.Message != "" && e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err)
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}

	return "error message not specified"
}

// GRPCStatus создает gRPC-ответ
//
//	Если error имеет данный метод, то он используется для вывода ошибки.
//	Позволяет скрывать ошибку Err из вывода gRPC (путем разделения Error() и GRPCStatus()).
//
//	Если Message не был указан, то выводит "Что-то пошло не так. Попробуйте еще раз."
func (e *apiError) GRPCStatus() *status.Status {
	const defaultMessage = "Что-то пошло не так. Попробуйте еще раз."

	if e.Message == "" {
		e.Message = defaultMessage
	}
	statusErr := status.New(e.Code, e.Message)
	var details []proto.Message
	if e.localizeMessage != nil {
		details = append(details, e.localizeMessage)
	}
	if len(e.Details) > 0 {
		details = append(details, e.Details...)
	}
	if len(details) > 0 {
		var err error
		statusErr, err = statusErr.WithDetails(details...)
		_ = err
	}

	return statusErr
}

func (e *apiError) WithDetails(details ...proto.Message) *apiError {
	e.Details = append(e.Details, details...)
	return e
}

func (e *apiError) WithLocalizedMessage(msg string) *apiError {
	if msg == "" {
		return e
	}
	e.localizeMessage = &errdetails.LocalizedMessage{
		Locale:  "ru-RU",
		Message: msg,
	}

	return e
}
