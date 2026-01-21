/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-21 09:58:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-21 09:58:15
 * @FilePath: \go-risk\store\store.go
 * @Description: 内存计数后端
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

const shardCount = 256

// Memory 内存计数后端，基于 256 分片降低锁竞争，单机/测试场景使用。
type Memory struct {
	shards [shardCount]*shard
	banMu  sync.RWMutex
	bans   map[string]core.BanEntry
}

// shard 单个分片计数器。
type shard struct {
	mu sync.Mutex
	m  map[string]*counter
}

// NewMemory 创建内存计数后端。
func NewMemory() *Memory {
	m := &Memory{bans: make(map[string]core.BanEntry)}
	for i := range m.shards {
		m.shards[i] = &shard{m: make(map[string]*counter)}
	}
	return m
}

// IncrBatch 批量递增计数。
func (m *Memory) IncrBatch(ctx context.Context, ops []core.CounterOp) error {
	for _, op := range ops {
		s := m.shardOf(op.Key)
		s.mu.Lock()
		c := s.m[op.Key]
		if c == nil {
			c = &counter{}
			s.m[op.Key] = c
		}
		c.incr(op.Delta, windowNano)
		s.mu.Unlock()
	}
	return nil
}

// WindowCount 返回窗口内计数。
func (m *Memory) WindowCount(ctx context.Context, key string, window time.Duration) (int64, error) {
	s := m.shardOf(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.m[key]
	if c == nil {
		return 0, nil
	}
	return c.read(window), nil
}

// LoadBanList 内存后端无持久化封禁，返回空。
func (m *Memory) LoadBanList(ctx context.Context) ([]core.BanEntry, error) {
	m.banMu.RLock()
	defer m.banMu.RUnlock()
	out := make([]core.BanEntry, 0, len(m.bans))
	for _, b := range m.bans {
		out = append(out, b)
	}
	return out, nil
}

func (m *Memory) shardOf(key string) *shard {
	return m.shards[fnv32(key)%shardCount]
}

func fnv32(key string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}
	return h
}