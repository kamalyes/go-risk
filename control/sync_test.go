/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-23 09:35:17
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-23 09:35:17
 * @FilePath: \go-risk\control\sync_test.go
 * @Description: 规则同步器单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package control

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

// fakeSyncerNotifier 记录发布并捕获订阅回调的通知器桩
type fakeSyncerNotifier struct {
	topic   string
	payload []byte
	handler func([]byte)
}

func (f *fakeSyncerNotifier) Publish(_ context.Context, topic string, payload []byte) error {
	f.topic = topic
	f.payload = payload
	return nil
}

func (f *fakeSyncerNotifier) Subscribe(_ context.Context, _ string, handler func([]byte)) error {
	f.handler = handler
	return nil
}

func TestSyncerApplyReloadsAndPublishes(t *testing.T) {
	mgr := New(core.RuleSnapshot{Version: "v1", Rules: []core.Rule{{ID: "a", Priority: 1, Enabled: true}}})
	n := &fakeSyncerNotifier{}
	s := NewSyncer(mgr, n, "")

	snapshot := core.RuleSnapshot{Version: "v2", Rules: []core.Rule{{ID: "b", Priority: 1, Enabled: true}}}
	assert.NoError(t, s.Apply(context.Background(), snapshot))

	assert.Equal(t, "b", mgr.Current().Rules[0].ID)
	assert.Equal(t, DefaultRuleSyncTopic, n.topic)

	var got core.RuleSnapshot
	assert.NoError(t, json.Unmarshal(n.payload, &got))
	assert.Equal(t, "v2", got.Version)
	assert.Equal(t, "b", got.Rules[0].ID)
}

func TestSyncerRunAppliesRemote(t *testing.T) {
	mgr := New(core.RuleSnapshot{Version: "v1", Rules: []core.Rule{{ID: "a", Priority: 1, Enabled: true}}})
	n := &fakeSyncerNotifier{}
	s := NewSyncer(mgr, n, "")

	assert.NoError(t, s.Run(context.Background()))
	assert.NotNil(t, n.handler)

	payload, _ := json.Marshal(core.RuleSnapshot{Version: "v3", Rules: []core.Rule{{ID: "c", Priority: 1, Enabled: true}}})
	n.handler(payload)

	assert.Equal(t, "c", mgr.Current().Rules[0].ID)
	assert.Empty(t, n.topic)
}

func TestSyncerApplyWithoutNotifier(t *testing.T) {
	mgr := New(core.RuleSnapshot{Version: "v1", Rules: nil})
	s := NewSyncer(mgr, nil, "")

	snapshot := core.RuleSnapshot{Version: "v2", Rules: []core.Rule{{ID: "z", Priority: 1, Enabled: true}}}
	assert.NoError(t, s.Apply(context.Background(), snapshot))
	assert.Equal(t, "z", mgr.Current().Rules[0].ID)
}
