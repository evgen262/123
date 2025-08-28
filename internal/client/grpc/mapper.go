package grpc

import (
	"time"

	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func NewSharedMapper(tu timeUtils.TimeUtils) *sharedMapper {
	return &sharedMapper{
		tu: tu,
	}
}

type sharedMapper struct {
	tu timeUtils.TimeUtils
}

func (m *sharedMapper) UUIDStringValue(u *uuid.UUID) *wrapperspb.StringValue {
	if u == nil {
		return nil
	}

	return &wrapperspb.StringValue{
		Value: u.String(),
	}
}

func (m *sharedMapper) StringValue(s *string) *wrapperspb.StringValue {
	if s == nil {
		return nil
	}

	return &wrapperspb.StringValue{
		Value: *s,
	}
}

func (m *sharedMapper) Int32Value(i *int) *wrapperspb.Int32Value {
	if i == nil {
		return nil
	}

	return &wrapperspb.Int32Value{
		Value: int32(*i),
	}
}

func (m *sharedMapper) Int64Value(i *int64) *wrapperspb.Int64Value {
	if i == nil {
		return nil
	}

	return &wrapperspb.Int64Value{
		Value: *i,
	}
}

func (m *sharedMapper) BoolValue(b *bool) *wrapperspb.BoolValue {
	if b == nil {
		return nil
	}

	return &wrapperspb.BoolValue{
		Value: *b,
	}
}

func (m *sharedMapper) BytesValue(b []byte) *wrapperspb.BytesValue {
	if b == nil {
		return nil
	}

	return &wrapperspb.BytesValue{
		Value: b,
	}
}

func (m *sharedMapper) TimeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	return m.tu.TimeToTimestamp(t)
}

func (m *sharedMapper) TimestampToTime(t *timestamppb.Timestamp) *time.Time {
	return m.tu.TimestampToTime(t)
}

func (m *sharedMapper) Int32SliceToInt(v []int32) []int {
	if len(v) == 0 {
		return nil
	}
	ints := make([]int, 0, len(v))
	for _, v := range v {
		ints = append(ints, int(v))
	}
	return ints
}

func (m *sharedMapper) IntSliceToInt32(v []int) []int32 {
	if len(v) == 0 {
		return nil
	}
	ints32 := make([]int32, 0, len(v))
	for _, v := range v {
		ints32 = append(ints32, int32(v))
	}
	return ints32
}

func (m *sharedMapper) StringSliceToUUID(v []string, skipZero bool) []uuid.UUID {
	if len(v) == 0 {
		return nil
	}
	uuids := make([]uuid.UUID, 0, len(v))
	for _, v := range v {
		id, err := uuid.Parse(v)
		if err != nil {
			id = uuid.Nil
		}
		if skipZero && id == uuid.Nil {
			continue
		}
		uuids = append(uuids, id)
	}
	return uuids
}

func (m *sharedMapper) StringSliceToPtrUUID(v []string, skipZero bool) []*uuid.UUID {
	if len(v) == 0 {
		return nil
	}
	uuids := make([]*uuid.UUID, 0, len(v))
	for _, v := range v {
		id, err := uuid.Parse(v)
		if err != nil {
			id = uuid.Nil
		}
		if skipZero && id == uuid.Nil {
			continue
		}
		uuids = append(uuids, &id)
	}
	return uuids
}

func (m *sharedMapper) StringToStringPtr(v string, skipZero bool) *string {
	if v == "" {
		if skipZero {
			return nil
		}
		return &v
	}

	return &v
}

func (m *sharedMapper) StringValueToStringPtr(v *wrapperspb.StringValue) *string {
	if v == nil {
		return nil
	}

	s := v.GetValue()
	return &s
}

func (m *sharedMapper) StringToUUIDPtr(v string, skipZero bool) *uuid.UUID {
	if v == "" {
		if skipZero {
			return nil
		}
		return &uuid.Nil
	}

	s, err := uuid.Parse(v)
	if err != nil {
		s = uuid.Nil
	}
	if skipZero && s == uuid.Nil {
		return nil
	}
	return &s
}

func (m *sharedMapper) StringValueToUUIDPtr(v *wrapperspb.StringValue, skipZero bool) *uuid.UUID {
	if v == nil {
		return nil
	}

	s, err := uuid.Parse(v.GetValue())
	if err != nil {
		s = uuid.Nil
	}
	if skipZero && s == uuid.Nil {
		return nil
	}
	return &s
}
