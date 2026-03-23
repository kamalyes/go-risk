/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-23 10:19:38
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-23 10:19:38
 * @FilePath: \go-risk\banlist\banlist.go
 * @Description: 本地封禁名单，读写分离维护封禁条目
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package banlist

import (
	"sync"
	"time"

	"github.com/kamalyes/go-risk/core"
)

// List 本地封禁名单，读路径无锁、写路径串行，供拦截 Slot 热路径查询
type List struct {
	mu      sync.RWMutex
	entries map[string]core.BanEntry
}

// New 创建空封禁名单
func New() *List {
	return &List{entries: make(map[string]core.BanEntry)}
}

// Apply 批量写入封禁条目，同键覆盖以刷新过期时间
func (l *List) Apply(entries []core.BanEntry) {
	l.mu.Lock()
	for _, e := range entries {
		l.entries[e.Key] = e
	}
	l.mu.Unlock()
}

// Remove 批量移除封禁条目，用于解禁
func (l *List) Remove(entries []core.BanEntry) {
	l.mu.Lock()
	for _, e := range entries {
		delete(l.entries, e.Key)
	}
	l.mu.Unlock()
}

// Lookup 查询封禁条目，过期条目即时剔除并返回未命中
func (l *List) Lookup(key string) (core.BanEntry, bool) {
	l.mu.RLock()
	e, ok := l.entries[key]
	l.mu.RUnlock()
	if !ok {
		return core.BanEntry{}, false
	}
	if !e.ExpireAt.IsZero() && time.Now().After(e.ExpireAt) {
		l.mu.Lock()
		delete(l.entries, key)
		l.mu.Unlock()
		return core.BanEntry{}, false
	}
	return e, true
}
