package kadry

import (
	"encoding/json"

	"go.uber.org/zap/zapcore"
)

type commonResponse struct {
	MessageType      string `json:"MessageType,omitempty"`
	RequestExecuted  bool   `json:"RequestExecuted"`
	RequestType      string `json:"RequestType,omitempty"`
	ErrorDescription string `json:"ErrorDescription,omitempty"`
}

type response struct {
	commonResponse
	Body struct {
		MobileApp mobileApp `json:"MobileApp,omitempty"`
	} `json:"ResponceBody,omitempty"`
}

func (r *response) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	enc.AddByteString("json", data)
	return nil
}
