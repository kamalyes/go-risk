/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-17 18:52:39
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-17 18:52:39
 * @FilePath: \go-risk\core\interfaces.go
 * @Description: 风控引擎、计数后端与通知器接口契约
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

import (
	"context"
	"time"
)

// Engine 风控引擎统一入口。适配层只面向该接口，核心逻辑与框架解耦。
type Engine interface {
	// Evaluate 同步判定（本地快路径）；跨实例聚合由引擎内部 Store 完成。
	Evaluate(ctx context.Context, rc *RiskContext) Decision

	// MarkResult 请求完成后回写结果，用于失败率、暴力破解等行为分析。
	MarkResult(ctx context.Context, rc *RiskContext, statusCode int, ok bool)

	// Close 释放订阅等资源。
	Close() error
}

// CounterOp 一次计数操作，支持批量异步攒批聚合。
type CounterOp struct {
	Key   string
	Delta int64
}

// CounterStore 分布式计数后端抽象（内存实现 + Redis 实现）。
type CounterStore interface {
	// IncrBatch 批量递增计数。
	IncrBatch(ctx context.Context, ops []CounterOp) error

	// WindowCount 统计窗口内计数。
	WindowCount(ctx context.Context, key string, window time.Duration) (int64, error)

	// LoadBanList 加载封禁名单（用于实例启动预热）。
	LoadBanList(ctx context.Context) ([]BanEntry, error)
}

// BanEntry 封禁条目。
type BanEntry struct {
	Key      string
	Scope    string
	ExpireAt time.Time
}

// Notifier 封禁/解禁广播抽象（内存实现 + NATS 实现）。
type Notifier interface {
	// Publish 发布事件。
	Publish(ctx context.Context, topic string, payload []byte) error

	// Subscribe 订阅事件，handler 在收到消息时回调。
	Subscribe(ctx context.Context, topic string, handler func([]byte)) error
}