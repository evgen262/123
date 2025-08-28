package dto

type OrderDirection string

const (
	OrderDirectionAsc  OrderDirection = "asc"
	OrderDirectionDesc OrderDirection = "desc"
)

type SortField int

const (
	SortFieldUnknown    SortField = 0
	SortFieldDateCreate SortField = 1
)
