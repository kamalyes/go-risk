/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-21 22:05:23
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-21 22:05:23
 * @FilePath: \go-risk\store\fallback.go
 * @Description: 计数后端降级兜底
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package store

import (
	"context"
	"time"

	"github.com/kamalyes/go-risk/core"
)

// Fallback 降级装饰器：主后端（如 Redis）失败时降级到备用后端（如内存），保证风控不中断
type Fallback struct {
	primary core.CounterStore
	backup  core.CounterStore
}

// NewFallback 创建降级后端
func NewFallback(primary, backup core.CounterStore) core.CounterStore {
	return &Fallback{primary: primary, backup: backup}
}

// IncrBatch 主后端失败时降级到备用后端
func (f *Fallback) IncrBatch(ctx context.Context, ops []core.CounterOp) error {
	if err := f.primary.IncrBatch(ctx, ops); err != nil {
		return f.backup.IncrBatch(ctx, ops)
	}
	return nil
}

// WindowCount 主后端失败时降级到备用后端
func (f *Fallback) WindowCount(ctx context.Context, key string, window time.Duration) (int64, error) {
	c, err := f.primary.WindowCount(ctx, key, window)
	if err != nil {
		return f.backup.WindowCount(ctx, key, window)
	}
	return c, nil
}

// LoadBanList 主后端失败时降级到备用后端
func (f *Fallback) LoadBanList(ctx context.Context) ([]core.BanEntry, error) {
	list, err := f.primary.LoadBanList(ctx)
	if err != nil {
		return f.backup.LoadBanList(ctx)
	}
	return list, nil
}