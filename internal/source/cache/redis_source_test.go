package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type redisSuite struct {
	suite.Suite

	redis *MockRedis

	ctx context.Context
}

func (rs *redisSuite) SetupTest() {
	ctrl := gomock.NewController(rs.T())

	rs.redis = NewMockRedis(ctrl)

	rs.ctx = context.TODO()
}

func (rs *redisSuite) Test_Get_Err() {
	c := NewCacheSource(rs.redis)

	wantErr := errors.New("some redis error")

	testCmd := &redis.StringCmd{}
	testCmd.SetErr(wantErr)

	rs.redis.EXPECT().Get(rs.ctx, "test_key").Return(testCmd)
	got, err := c.Get(rs.ctx, "test_key")

	rs.EqualError(err, wantErr.Error())
	rs.Empty(got)
}

func (rs *redisSuite) Test_Get_Ok() {
	c := NewCacheSource(rs.redis)

	want := "some_test_value"

	testCmd := &redis.StringCmd{}
	testCmd.SetVal(want)

	rs.redis.EXPECT().Get(rs.ctx, "test_key").Return(testCmd)
	got, err := c.Get(rs.ctx, "test_key")

	rs.Equal(want, got)
	rs.NoError(err)
}

func (rs *redisSuite) Test_SetEx_Err() {
	c := NewCacheSource(rs.redis)

	wantErr := errors.New("some redis error")

	testCmd := &redis.StatusCmd{}
	testCmd.SetErr(wantErr)

	rs.redis.EXPECT().SetEx(rs.ctx, "test_key", "some_test_value", 5*time.Second).Return(testCmd)
	err := c.SetEx(rs.ctx, "test_key", "some_test_value", 5*time.Second)

	rs.EqualError(err, wantErr.Error())
}

func (rs *redisSuite) Test_SetEx_Ok() {
	c := NewCacheSource(rs.redis)

	testCmd := &redis.StatusCmd{}

	rs.redis.EXPECT().SetEx(rs.ctx, "test_key", "some_test_value", 5*time.Second).Return(testCmd)
	err := c.SetEx(rs.ctx, "test_key", "some_test_value", 5*time.Second)

	rs.NoError(err)
}

func (rs *redisSuite) Test_Exists_Err() {
	c := NewCacheSource(rs.redis)

	wantErr := errors.New("some redis error")

	testCmd := &redis.IntCmd{}
	testCmd.SetErr(wantErr)

	rs.redis.EXPECT().Exists(rs.ctx, "test_key_1", "test_key_2").Return(testCmd)
	got, err := c.Exists(rs.ctx, "test_key_1", "test_key_2")

	rs.Empty(got)
	rs.EqualError(err, wantErr.Error())
}

func (rs *redisSuite) Test_Exists_False() {
	c := NewCacheSource(rs.redis)

	testCmd := &redis.IntCmd{}
	testCmd.SetVal(1)

	rs.redis.EXPECT().Exists(rs.ctx, "test_key_1", "test_key_2").Return(testCmd)
	got, err := c.Exists(rs.ctx, "test_key_1", "test_key_2")

	rs.Equal(false, got)
	rs.NoError(err)
}

func (rs *redisSuite) Test_Exists_True() {
	c := NewCacheSource(rs.redis)

	testCmd := &redis.IntCmd{}
	testCmd.SetVal(2)

	rs.redis.EXPECT().Exists(rs.ctx, "test_key_1", "test_key_2").Return(testCmd)
	got, err := c.Exists(rs.ctx, "test_key_1", "test_key_2")

	rs.Equal(true, got)
	rs.NoError(err)
}

func (rs *redisSuite) Test_Del_Err() {
	c := NewCacheSource(rs.redis)

	wantErr := errors.New("some redis error")

	testCmd := &redis.IntCmd{}
	testCmd.SetErr(wantErr)

	rs.redis.EXPECT().Del(rs.ctx, "test_key_1", "test_key_2").Return(testCmd)
	err := c.Delete(rs.ctx, "test_key_1", "test_key_2")

	rs.EqualError(err, wantErr.Error())
}

func (rs *redisSuite) Test_Del_False() {
	c := NewCacheSource(rs.redis)

	testCmd := &redis.IntCmd{}
	testCmd.SetVal(1)

	rs.redis.EXPECT().Del(rs.ctx, "test_key_1", "test_key_2").Return(testCmd)
	err := c.Delete(rs.ctx, "test_key_1", "test_key_2")

	rs.NoError(err)
}

func (rs *redisSuite) Test_Del_True() {
	c := NewCacheSource(rs.redis)

	testCmd := &redis.IntCmd{}
	testCmd.SetVal(2)

	rs.redis.EXPECT().Del(rs.ctx, "test_key_1", "test_key_2").Return(testCmd)
	err := c.Delete(rs.ctx, "test_key_1", "test_key_2")

	rs.NoError(err)
}

func TestSync(t *testing.T) {
	suite.Run(t, &redisSuite{})
}
