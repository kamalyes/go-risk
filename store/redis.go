/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-22 10:12:31
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-22 10:12:31
 * @FilePath: \go-risk\store\redis.go
 * @Description: Redis 分布式计数后端
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package store

import (
	"context"
	"time"

	"github.com/kamalyes/go-risk/core"
)

// RedisClient Redis 客户端最小契约，核心库不直接依赖 go-redis 等具体 SDK
type RedisClient interface {
	// IncrBy 原子递增并返回最新值
	IncrBy(ctx context.Context, key string, delta int64) (int64, error)

	// Expire 设置键过期时间
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// Get 读取计数，键不存在时返回 0
	Get(ctx context.Context, key string) (int64, error)
}

// Redis 基于 Redis 的分布式计数后端，跨实例共享固定窗口计数
type Redis struct {
	client RedisClient
	window time.Duration
}

// NewRedis 创建 Redis 计数后端，window 为计数窗口（<=0 时默认 1 秒）
func NewRedis(client RedisClient, window time.Duration) *Redis {
	if window <= 0 {
		window = time.Second
	}
	return &Redis{client: client, window: window}
}

// IncrBatch 批量递增计数，首次写入时设置 TTL 避免热键反复续期
func (r *Redis) IncrBatch(ctx context.Context, ops []core.CounterOp) error {
	for _, op := range ops {
		val, err := r.client.IncrBy(ctx, op.Key, op.Delta)
		if err != nil {
			return err
		}
		if val == op.Delta {
			if err := r.client.Expire(ctx, op.Key, r.window); err != nil {
				return err
			}
		}
	}
	return nil
}

// WindowCount 返回窗口内计数，过期由 Redis TTL 自行回收
func (r *Redis) WindowCount(ctx context.Context, key string, window time.Duration) (int64, error) {
	return r.client.Get(ctx, key)
}

// LoadBanList 封禁名单经 NATS 广播下发，Redis 不持久化封禁条目
func (r *Redis) LoadBanList(ctx context.Context) ([]core.BanEntry, error) {
	return nil, nil
}
