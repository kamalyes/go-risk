/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-21 09:58:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-21 17:22:31
 * @FilePath: \go-risk\store\memory.go
 * @Description: 分片固定窗口计数器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package store

import (
	"sync/atomic"
	"time"
)

// windowNano 默认固定窗口 1 秒。
const windowNano = int64(time.Second)

// counter 固定窗口计数器（无锁原子操作）。
type counter struct {
	count         int64 // 计数
	resetTimeNano int64 // 窗口重置时间
}

// incr 递增计数；窗口回绕时通过 CAS 保证仅一个 goroutine 重置，
// 避免并发下重复清零导致计数丢失。
func (c *counter) incr(delta int64, window int64) {
	for {
		now := time.Now().UnixNano()
		rt := atomic.LoadInt64(&c.resetTimeNano)
		if now > rt {
			if atomic.CompareAndSwapInt64(&c.resetTimeNano, rt, now+window) {
				atomic.StoreInt64(&c.count, 0)
			}
			continue
		}
		break
	}
	atomic.AddInt64(&c.count, delta)
}

// read 返回当前计数（窗口过期则返回 0）。
func (c *counter) read(window time.Duration) int64 {
	now := time.Now().UnixNano()
	rt := atomic.LoadInt64(&c.resetTimeNano)
	if now > rt {
		return 0
	}
	return atomic.LoadInt64(&c.count)
}
