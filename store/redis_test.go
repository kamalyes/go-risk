/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-23 20:15:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-23 20:15:33
 * @FilePath: \go-risk\store\redis_test.go
 * @Description: Redis 计数后端单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

// fakeRedis 内存版 Redis 客户端，用于验证 Redis 后端的编排逻辑
type fakeRedis struct {
	vals      map[string]int64
	incrErr   error
	expireErr error
	expireKey string
	expireTTL time.Duration
}

func newFakeRedis() *fakeRedis {
	return &fakeRedis{vals: make(map[string]int64)}
}

func (f *fakeRedis) IncrBy(_ context.Context, key string, delta int64) (int64, error) {
	if f.incrErr != nil {
		return 0, f.incrErr
	}
	f.vals[key] += delta
	return f.vals[key], nil
}

func (f *fakeRedis) Expire(_ context.Context, key string, ttl time.Duration) error {
	if f.expireErr != nil {
		return f.expireErr
	}
	f.expireKey = key
	f.expireTTL = ttl
	return nil
}

func (f *fakeRedis) Get(_ context.Context, key string) (int64, error) {
	return f.vals[key], nil
}

func TestRedisDefaultWindow(t *testing.T) {
	f := newFakeRedis()
	r := NewRedis(f, 0)
	err := r.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:a", Delta: 1}})
	assert.NoError(t, err)
	assert.Equal(t, time.Second, f.expireTTL)
}

func TestRedisIncrBatchSetTTLOnFirstWrite(t *testing.T) {
	f := newFakeRedis()
	r := NewRedis(f, time.Second)

	err := r.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:b", Delta: 1}})
	assert.NoError(t, err)
	assert.Equal(t, "ip:b", f.expireKey)
	assert.Equal(t, time.Second, f.expireTTL)

	// 再次写入不再重复续期
	f.expireKey = ""
	err = r.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:b", Delta: 1}})
	assert.NoError(t, err)
	assert.Equal(t, "", f.expireKey)
	assert.Equal(t, int64(2), f.vals["ip:b"])
}

func TestRedisIncrBatchIncrError(t *testing.T) {
	f := newFakeRedis()
	f.incrErr = errors.New("incr down")
	r := NewRedis(f, time.Second)
	err := r.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:c", Delta: 1}})
	assert.Error(t, err)
}

func TestRedisIncrBatchExpireError(t *testing.T) {
	f := newFakeRedis()
	f.expireErr = errors.New("expire down")
	r := NewRedis(f, time.Second)
	err := r.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:d", Delta: 1}})
	assert.Error(t, err)
}

func TestRedisWindowCount(t *testing.T) {
	f := newFakeRedis()
	r := NewRedis(f, time.Second)
	_ = r.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:e", Delta: 3}})

	c, err := r.WindowCount(context.Background(), "ip:e", time.Second)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), c)
}

func TestRedisLoadBanListEmpty(t *testing.T) {
	f := newFakeRedis()
	r := NewRedis(f, time.Second)
	list, err := r.LoadBanList(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, list)
}