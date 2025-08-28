package http

type MiddlewareOption struct {
	Name  string
	Value interface{}
}

func (mo MiddlewareOption) GetName() string {
	return mo.Name
}

func (mo MiddlewareOption) GetValue() interface{} {
	return mo.Value
}

func (mo MiddlewareOption) String() string {
	val, ok := mo.Value.(string)
	if !ok {
		return ""
	}

	return val
}

type MiddlewareOptions struct {
	opts []*MiddlewareOption
}

func (o *MiddlewareOptions) Add(option *MiddlewareOption) {
	o.opts = append(o.opts, option)
}

func (o *MiddlewareOptions) Get(optName string) *MiddlewareOption {
	for _, option := range o.opts {
		if option.Name == optName {
			return option
		}
	}

	return nil
}

func (o *MiddlewareOptions) Filter(optNames ...string) []*MiddlewareOption {
	opts := make([]*MiddlewareOption, 0, len(optNames))
	for _, optName := range optNames {
		opt := o.Get(optName)
		if opt != nil {
			opts = append(opts, opt)
		}
	}

	return opts
}
