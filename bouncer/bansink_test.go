/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-22 15:27:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-22 15:27:33
 * @FilePath: \go-risk\bouncer\bansink_test.go
 * @Description: 通知器封禁落地单测（广播主题与载荷序列化）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/kamalyes/go-risk/notifier"
	"github.com/stretchr/testify/assert"
)

// fakeNotifier 记录最近一次发布主题与载荷的通知器桩
type fakeNotifier struct {
	topic   string
	payload []byte
}

func (f *fakeNotifier) Publish(_ context.Context, topic string, payload []byte) error {
	f.topic = topic
	f.payload = payload
	return nil
}

func (f *fakeNotifier) Subscribe(_ context.Context, _ string, _ func([]byte)) error {
	return nil
}

func TestNotifierSinkBan(t *testing.T) {
	n := &fakeNotifier{}
	sink := NewNotifierSink(n)
	entries := []core.BanEntry{{Key: "ip:1.2.3.11", Scope: "ip", ExpireAt: time.Now()}}

	err := sink.Ban(context.Background(), entries)
	assert.NoError(t, err)
	assert.Equal(t, notifier.TopicBan, n.topic)

	var got []core.BanEntry
	assert.NoError(t, json.Unmarshal(n.payload, &got))
	assert.Len(t, got, 1)
	assert.Equal(t, "ip:1.2.3.11", got[0].Key)
}

func TestNotifierSinkUnban(t *testing.T) {
	n := &fakeNotifier{}
	sink := NewNotifierSink(n)

	err := sink.Unban(context.Background(), []core.BanEntry{{Key: "ip:1.2.3.11", Scope: "ip"}})
	assert.NoError(t, err)
	assert.Equal(t, notifier.TopicUnban, n.topic)
}

func TestNotifierSinkNilNotifier(t *testing.T) {
	sink := NewNotifierSink(nil)

	assert.NoError(t, sink.Ban(context.Background(), nil))
	assert.NoError(t, sink.Unban(context.Background(), nil))
}
