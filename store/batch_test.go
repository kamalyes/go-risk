/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-25 19:37:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-25 19:37:16
 * @FilePath: \go-risk\store\batch_test.go
 * @Description: 攒批计数后端单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package store

import (
	"context"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestBatchingFlushOnThreshold(t *testing.T) {
	m := NewMemory()
	b := NewBatching(m, 2, time.Hour)
	defer b.Close()

	err := b.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:a", Delta: 1}})
	assert.NoError(t, err)
	c, _ := m.WindowCount(context.Background(), "ip:a", time.Second)
	assert.Equal(t, int64(0), c)

	err = b.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:a", Delta: 1}})
	assert.NoError(t, err)
	c, _ = m.WindowCount(context.Background(), "ip:a", time.Second)
	assert.Equal(t, int64(2), c)
}

func TestBatchingFlushEmpty(t *testing.T) {
	m := NewMemory()
	b := NewBatching(m, 2, time.Hour)
	err := b.Flush(context.Background())
	assert.NoError(t, err)
	b.Close()
}

func TestBatchingClose(t *testing.T) {
	m := NewMemory()
	b := NewBatching(m, 128, time.Hour)
	err := b.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:b", Delta: 1}})
	assert.NoError(t, err)
	err = b.Close()
	assert.NoError(t, err)
	c, _ := m.WindowCount(context.Background(), "ip:b", time.Second)
	assert.Equal(t, int64(1), c)
}

func TestBatchingWindowCountPassthrough(t *testing.T) {
	m := NewMemory()
	b := NewBatching(m, 2, time.Hour)
	defer b.Close()
	_ = m.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:c", Delta: 5}})

	c, err := b.WindowCount(context.Background(), "ip:c", time.Second)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), c)
}

func TestBatchingLoadBanListPassthrough(t *testing.T) {
	m := NewMemory()
	b := NewBatching(m, 2, time.Hour)
	defer b.Close()
	list, err := b.LoadBanList(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, list)
}

func TestBatchingDefaultParams(t *testing.T) {
	m := NewMemory()
	b := NewBatching(m, 0, 0)
	defer b.Close()

	ops := make([]core.CounterOp, 200)
	for i := range ops {
		ops[i] = core.CounterOp{Key: "ip:e", Delta: 1}
	}
	err := b.IncrBatch(context.Background(), ops)
	assert.NoError(t, err)
	c, _ := m.WindowCount(context.Background(), "ip:e", time.Second)
	assert.Equal(t, int64(200), c)
}

func TestBatchingTickerFlush(t *testing.T) {
	m := NewMemory()
	b := NewBatching(m, 128, 10*time.Millisecond)
	defer b.Close()
	err := b.IncrBatch(context.Background(), []core.CounterOp{{Key: "ip:f", Delta: 1}})
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	c, _ := m.WindowCount(context.Background(), "ip:f", time.Second)
	assert.Equal(t, int64(1), c)
}