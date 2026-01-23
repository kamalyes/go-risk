/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 09:31:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 09:31:22
 * @FilePath: \go-risk\engine\engine.go
 * @Description: 风控引擎装配与统一入口
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package engine

import (
	"context"

	"github.com/kamalyes/go-risk/core"
	"github.com/kamalyes/go-risk/notifier"
	"github.com/kamalyes/go-risk/store"
)

// Engine 风控引擎，装配 core 契约与各功能模块。
type Engine struct {
	cfg      *core.Config
	store    core.CounterStore
	notifier core.Notifier
	slots    []core.Slot
	chain    *slotChain
}

// Option 引擎配置项。
type Option func(*Engine)

// New 创建引擎，默认内存后端 + 空 Slot 链。
func New(opts ...Option) *Engine {
	e := &Engine{
		cfg:      core.DefaultConfig(),
		store:    store.NewMemory(),
		notifier: notifier.NewMemory(),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// WithConfig 覆盖默认配置。
func WithConfig(cfg *core.Config) Option {
	return func(e *Engine) {
		if cfg != nil {
			e.cfg = cfg
		}
	}
}

// WithStore 注入计数后端。
func WithStore(s core.CounterStore) Option {
	return func(e *Engine) {
		if s != nil {
			e.store = s
		}
	}
}

// WithNotifier 注入通知器。
func WithNotifier(n core.Notifier) Option {
	return func(e *Engine) {
		if n != nil {
			e.notifier = n
		}
	}
}

// WithSlots 追加 Slot。
func WithSlots(slots ...core.Slot) Option {
	return func(e *Engine) {
		e.slots = append(e.slots, slots...)
	}
}

// Evaluate 同步决策。
func (e *Engine) Evaluate(ctx context.Context, rc *core.RiskContext) core.Decision {
	return e.chain.run(ctx, rc)
}

// MarkResult 请求完成后回写结果，行为分析在后续里程碑落地。
func (e *Engine) MarkResult(ctx context.Context, rc *core.RiskContext, statusCode int, ok bool) {}

// Close 释放资源。
func (e *Engine) Close() error { return nil }