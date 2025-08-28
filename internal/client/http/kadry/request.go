package kadry

import (
	"encoding/json"

	"go.uber.org/zap/zapcore"
)

type attributes struct {
	Exclude []AttributeName `json:"Exclude,omitempty"`
	Include []AttributeName `json:"Include,omitempty"`
}

type attributeList struct {
	MobileApp attributes `json:"MobileApp,omitempty"`
}

type request struct {
	PersonIDArray []string       `json:"PersonIDArray"`
	AttributeList *attributeList `json:"AttributeList,omitempty"`
}

func (r *request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	enc.AddByteString("json", data)
	return nil
}

func newRequest(persons []string, attrs ...AttributeName) *request {
	req := &request{
		PersonIDArray: persons,
	}

	if len(attrs) > 0 {
		req.AttributeList = new(attributeList)

		for _, f := range attrs {
			req.AttributeList.MobileApp.Include = append(req.AttributeList.MobileApp.Include, f)
		}
	}

	return req
}
