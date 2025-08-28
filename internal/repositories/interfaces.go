package repositories

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//go:generate mockgen -source=interfaces.go -destination=repositories_mock.go -package=repositories

type SharedMapper interface {
	UUIDStringValue(u *uuid.UUID) *wrapperspb.StringValue
	StringValue(s *string) *wrapperspb.StringValue
	Int32Value(i *int) *wrapperspb.Int32Value
	Int64Value(i *int64) *wrapperspb.Int64Value
	BoolValue(b *bool) *wrapperspb.BoolValue
	BytesValue(b []byte) *wrapperspb.BytesValue
	TimeToTimestamp(t *time.Time) *timestamppb.Timestamp
	TimestampToTime(t *timestamppb.Timestamp) *time.Time
	Int32SliceToInt(v []int32) []int
	IntSliceToInt32(v []int) []int32
	StringSliceToUUID(v []string, skipZero bool) []uuid.UUID
	StringSliceToPtrUUID(v []string, skipZero bool) []*uuid.UUID
	StringValueToStringPtr(v *wrapperspb.StringValue) *string
	StringValueToUUIDPtr(v *wrapperspb.StringValue, skipZero bool) *uuid.UUID
	StringToStringPtr(v string, skipZero bool) *string
	StringToUUIDPtr(v string, skipZero bool) *uuid.UUID
}
