/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-22 10:12:31
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-22 15:38:05
 * @FilePath: \go-risk\notifier\notifier.go
 * @Description: 内存广播通知器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package notifier

import "sync"

// 封禁/解禁主题，统一在此维护避免散落魔法字符串。
const (
	TopicBan   = "risk.ban"
	TopicUnban = "risk.unban"
)

// Memory 内存通知器，进程内广播封禁/解禁事件，单机场景使用。
type Memory struct {
	mu   sync.RWMutex
	subs map[string][]func([]byte)
}

// NewMemory 创建内存通知器。
func NewMemory() *Memory {
	return &Memory{subs: make(map[string][]func([]byte))}
}
