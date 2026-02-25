/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-25 19:37:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-25 19:37:16
 * @FilePath: \go-risk\store\fallback_test.go
 * @Description: 降级计数后端单元测试
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

// errStore 恒失败的计数后端，用于验证降级路径
type errStore struct{ err error }

func (e *errStore) IncrBatch(context.Context, []core.CounterOp) error { return e.err }
func (e *errStore) WindowCount(context.Context, string, time.Duration) (int64, error) {
	return 0, e.err
}
func (e *errStore) LoadBanList(context.Context) ([]core.BanEntry, error) { return nil, e.err }

func TestFallbackIncrBatchPrimaryOK(t *testing.T) {
	primary := NewMemory()
	backup := NewMemory()
	f := NewFallback(primary, backup)
	err := f.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:a", Delta: 1}})
	assert.NoError(t, err)
	c, _ := primary.WindowCount(context.Background(), "ip:a", time.Second)
	assert.Equal(t, int64(1), c)
	c, _ = backup.WindowCount(context.Background(), "ip:a", time.Second)
	assert.Equal(t, int64(0), c)
}

func TestFallbackIncrBatchDegrade(t *testing.T) {
	primary := &errStore{err: errors.New("down")}
	backup := NewMemory()
	f := NewFallback(primary, backup)
	err := f.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:b", Delta: 1}})
	assert.NoError(t, err)
	c, _ := backup.WindowCount(context.Background(), "ip:b", time.Second)
	assert.Equal(t, int64(1), c)
}

func TestFallbackWindowCountPrimaryOK(t *testing.T) {
	primary := NewMemory()
	backup := NewMemory()
	f := NewFallback(primary, backup)
	_ = primary.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:c", Delta: 9}})

	c, err := f.WindowCount(context.Background(), "ip:c", time.Second)
	assert.NoError(t, err)
	assert.Equal(t, int64(9), c)
}

func TestFallbackWindowCountDegrade(t *testing.T) {
	primary := &errStore{err: errors.New("down")}
	backup := NewMemory()
	f := NewFallback(primary, backup)
	_ = backup.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:d", Delta: 7}})

	c, err := f.WindowCount(context.Background(), "ip:d", time.Second)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), c)
}

func TestFallbackLoadBanListPrimaryOK(t *testing.T) {
	primary := NewMemory()
	backup := NewMemory()
	f := NewFallback(primary, backup)
	list, err := f.LoadBanList(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, list)
}

func TestFallbackLoadBanListDegrade(t *testing.T) {
	primary := &errStore{err: errors.New("down")}
	backup := NewMemory()
	f := NewFallback(primary, backup)
	list, err := f.LoadBanList(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, list)
}