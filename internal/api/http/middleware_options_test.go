package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareOption_GetName(t *testing.T) {
	type args struct {
		Name  string
		Value interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test GetName method",
			args: args{
				Name:  "Test Name Value",
				Value: nil,
			},
			want: "Test Name Value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mo := MiddlewareOption{
				Name:  tt.args.Name,
				Value: tt.args.Value,
			}
			assert.Equalf(t, tt.want, mo.GetName(), "GetName()")
		})
	}
}

func TestMiddlewareOption_GetValue(t *testing.T) {
	type args struct {
		Name  string
		Value interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "if value type of string",
			args: args{
				Name:  "Test Name",
				Value: "Test String Value",
			},
			want: "Test String Value",
		},
		{
			name: "if value not type of string",
			args: args{
				Name:  "Test Name",
				Value: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mo := MiddlewareOption{
				Name:  tt.args.Name,
				Value: tt.args.Value,
			}
			assert.Equalf(t, tt.want, mo.GetValue(), "GetValue()")
		})
	}
}

func TestMiddlewareOption_String(t *testing.T) {
	type args struct {
		Name  string
		Value interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "if value type of string",
			args: args{
				Name:  "Test Name",
				Value: "Test String Value",
			},
			want: "Test String Value",
		},
		{
			name: "if value not type of string",
			args: args{
				Name:  "Test Name",
				Value: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mo := MiddlewareOption{
				Name:  tt.args.Name,
				Value: tt.args.Value,
			}
			assert.Equalf(t, tt.want, mo.String(), "String()")
		})
	}
}

func TestMiddlewareOptions_Add(t *testing.T) {
	type args struct {
		option  *MiddlewareOption
		options []*MiddlewareOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "add",
			args: args{
				option: &MiddlewareOption{
					Name:  "TestName",
					Value: "TestValue",
				},
				options: []*MiddlewareOption{
					{
						Name:  "TestOption",
						Value: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &MiddlewareOptions{
				opts: tt.args.options,
			}
			o.Add(tt.args.option)
		})
	}
}

func TestMiddlewareOptions_Filter(t *testing.T) {
	type args struct {
		optNames []string
		options  []*MiddlewareOption
	}
	tests := []struct {
		name string
		args args
		want []*MiddlewareOption
	}{
		{
			name: "if contains all required options",
			args: args{
				optNames: []string{
					"TestOption1",
					"TestOption2",
					"TestOption3",
				},
				options: []*MiddlewareOption{
					{
						Name:  "TestOption1",
						Value: "TestValue1",
					},
					{
						Name:  "TestOption2",
						Value: nil,
					},
					{
						Name:  "TestOption3",
						Value: "TestValue3",
					},
				},
			},
			want: []*MiddlewareOption{
				{
					Name:  "TestOption1",
					Value: "TestValue1",
				},
				{
					Name:  "TestOption2",
					Value: nil,
				},
				{
					Name:  "TestOption3",
					Value: "TestValue3",
				},
			},
		},
		{
			name: "if contains the required options",
			args: args{
				optNames: []string{
					"TestOption1",
					"TestOption3",
				},
				options: []*MiddlewareOption{
					{
						Name:  "TestOption1",
						Value: "TestValue1",
					},
					{
						Name:  "TestOption2",
						Value: nil,
					},
					{
						Name:  "TestOption3",
						Value: "TestValue3",
					},
				},
			},
			want: []*MiddlewareOption{
				{
					Name:  "TestOption1",
					Value: "TestValue1",
				},
				{
					Name:  "TestOption3",
					Value: "TestValue3",
				},
			},
		},
		{
			name: "if does not contain the required options",
			args: args{
				optNames: []string{
					"TestOption3",
					"TestOption4",
				},
				options: []*MiddlewareOption{
					{
						Name:  "TestOption1",
						Value: nil,
					},
					{
						Name:  "TestOption2",
						Value: nil,
					},
				},
			},
			want: []*MiddlewareOption{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &MiddlewareOptions{
				opts: tt.args.options,
			}
			assert.Equalf(t, tt.want, o.Filter(tt.args.optNames...), "Filter(%v)", tt.args.optNames)
		})
	}
}

func TestMiddlewareOptions_Get(t *testing.T) {
	type args struct {
		optName string
		options []*MiddlewareOption
	}
	tests := []struct {
		name string
		args args
		want *MiddlewareOption
	}{
		{
			name: "if contains the required option",
			args: args{
				optName: "TestName",
				options: []*MiddlewareOption{
					{
						Name:  "TestName",
						Value: nil,
					},
				},
			},
			want: &MiddlewareOption{
				Name:  "TestName",
				Value: nil,
			},
		},
		{
			name: "if does not contain the required option",
			args: args{
				optName: "TestName",
				options: []*MiddlewareOption{
					{
						Name:  "SomeName",
						Value: nil,
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &MiddlewareOptions{
				opts: tt.args.options,
			}
			assert.Equalf(t, tt.want, o.Get(tt.args.optName), "Get(%v)", tt.args.optName)
		})
	}
}
