/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-20 16:05:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-20 16:05:53
 * @FilePath: \go-risk\bouncer\bansink.go
 * @Description: 基于通知器的封禁落地实现
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"context"
	"encoding/json"

	"github.com/kamalyes/go-risk/core"
	"github.com/kamalyes/go-risk/notifier"
)

// NotifierSink 基于 core.Notifier 的封禁落地，把封禁或解禁事件广播到集群内其余实例
type NotifierSink struct {
	n core.Notifier
}

// NewNotifierSink 创建通知器封禁落地
func NewNotifierSink(n core.Notifier) *NotifierSink {
	return &NotifierSink{n: n}
}

// Ban 批量广播封禁事件
func (s *NotifierSink) Ban(ctx context.Context, entries []core.BanEntry) error {
	return s.broadcast(ctx, notifier.TopicBan, entries)
}

// Unban 批量广播解禁事件
func (s *NotifierSink) Unban(ctx context.Context, entries []core.BanEntry) error {
	return s.broadcast(ctx, notifier.TopicUnban, entries)
}

// broadcast 序列化整批条目后发布，整批提交避免逐条订阅竞态
func (s *NotifierSink) broadcast(ctx context.Context, topic string, entries []core.BanEntry) error {
	if s.n == nil {
		return nil
	}
	payload, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	return s.n.Publish(ctx, topic, payload)
}
