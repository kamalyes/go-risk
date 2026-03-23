/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-23 09:28:31
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-23 09:28:31
 * @FilePath: \go-risk\control\sync.go
 * @Description: 控制面规则同步器，跨实例广播快照并订阅热更新
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package control

import (
	"context"
	"encoding/json"

	"github.com/kamalyes/go-risk/core"
)

// Syncer 规则同步器，把控制面变更广播到集群其余实例并订阅远程快照热更新
type Syncer struct {
	mgr   *Manager
	n     core.Notifier
	topic string
}

// NewSyncer 创建规则同步器，topic 为空时采用默认同步主题
func NewSyncer(mgr *Manager, n core.Notifier, topic string) *Syncer {
	if topic == "" {
		topic = DefaultRuleSyncTopic
	}
	return &Syncer{mgr: mgr, n: n, topic: topic}
}

// Apply 本地热切换快照并广播到集群，保证分布式下规则变更单一生效且跨节点一致
func (s *Syncer) Apply(ctx context.Context, snapshot core.RuleSnapshot) error {
	s.mgr.Reload(snapshot)
	if s.n == nil {
		return nil
	}
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return s.n.Publish(ctx, s.topic, payload)
}

// Run 订阅远程规则快照，收到广播后仅本地 Reload 避免循环回发
func (s *Syncer) Run(ctx context.Context) error {
	if s.n == nil {
		return nil
	}
	return s.n.Subscribe(ctx, s.topic, func(payload []byte) {
		var snapshot core.RuleSnapshot
		if err := json.Unmarshal(payload, &snapshot); err != nil {
			return
		}
		s.mgr.Reload(snapshot)
	})
}
