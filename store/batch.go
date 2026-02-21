/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-21 20:31:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-21 20:31:07
 * @FilePath: \go-risk\store\batch.go
 * @Description: 异步攒批计数后端
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package store

import (
	"context"
	"sync"
	"time"

	"github.com/kamalyes/go-risk/core"
)

// Batching 攒批装饰器：将 IncrBatch 合并为批量写入，降低远端 RTT（如 Redis Pipeline）
// 达到批量阈值立即 flush，否则按间隔定时 flush
type Batching struct {
	inner      core.CounterStore
	mu         sync.Mutex
	ops        []core.CounterOp
	batchSize  int
	interval   time.Duration
	stop       chan struct{}
	once       sync.Once
}

// NewBatching 创建攒批后端
func NewBatching(inner core.CounterStore, batchSize int, interval time.Duration) *Batching {
	if batchSize <= 0 {
		batchSize = 128
	}
	if interval <= 0 {
		interval = time.Millisecond * 50
	}
	b := &Batching{inner: inner, batchSize: batchSize, interval: interval, stop: make(chan struct{})}
	go b.loop()
	return b
}

// IncrBatch 累积计数，达到阈值时立即 flush
func (b *Batching) IncrBatch(ctx context.Context, ops []core.CounterOp) error {
	b.mu.Lock()
	b.ops = append(b.ops, ops...)
	flush := len(b.ops) >= b.batchSize
	b.mu.Unlock()
	if flush {
		return b.Flush(ctx)
	}
	return nil
}

// Flush 立即将累积的计数写入底层后端
func (b *Batching) Flush(ctx context.Context) error {
	b.mu.Lock()
	ops := b.ops
	b.ops = nil
	b.mu.Unlock()
	if len(ops) == 0 {
		return nil
	}
	return b.inner.IncrBatch(ctx, ops)
}

// WindowCount 透传到底层后端
func (b *Batching) WindowCount(ctx context.Context, key string, window time.Duration) (int64, error) {
	return b.inner.WindowCount(ctx, key, window)
}

// LoadBanList 透传到底层后端
func (b *Batching) LoadBanList(ctx context.Context) ([]core.BanEntry, error) {
	return b.inner.LoadBanList(ctx)
}

// Close 停止定时器并 flush 剩余计数
func (b *Batching) Close() error {
	b.once.Do(func() { close(b.stop) })
	return b.Flush(context.Background())
}

func (b *Batching) loop() {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_ = b.Flush(context.Background())
		case <-b.stop:
			return
		}
	}
}