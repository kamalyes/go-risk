/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-23 10:52:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-23 10:52:26
 * @FilePath: \go-risk\banlist\banlist_test.go
 * @Description: 封禁名单、拦截 Slot 与广播消费端单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package banlist

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/kamalyes/go-risk/notifier"
	"github.com/stretchr/testify/assert"
)

func TestListApplyLookupAndRemove(t *testing.T) {
	list := New()
	now := time.Now()
	list.Apply([]core.BanEntry{
		{Key: "ip:10.0.0.12", Scope: "ip", ExpireAt: now.Add(time.Hour)},
		{Key: "fingerprint:f7", Scope: "fingerprint", ExpireAt: now.Add(30 * time.Minute)},
		{Key: "tenant:t8", Scope: "tenant", ExpireAt: time.Time{}},
	})

	e, ok := list.Lookup("ip:10.0.0.12")
	assert.True(t, ok)
	assert.Equal(t, "ip", e.Scope)

	// 永久封禁（零过期时间）命中
	e, ok = list.Lookup("tenant:t8")
	assert.True(t, ok)
	assert.Equal(t, "tenant", e.Scope)

	// 未命中
	_, ok = list.Lookup("ip:10.0.0.99")
	assert.False(t, ok)

	// 解禁后未命中
	list.Remove([]core.BanEntry{{Key: "ip:10.0.0.12", Scope: "ip"}})
	_, ok = list.Lookup("ip:10.0.0.12")
	assert.False(t, ok)
}

func TestListExpiry(t *testing.T) {
	list := New()
	list.Apply([]core.BanEntry{{Key: "ip:10.0.0.12", Scope: "ip", ExpireAt: time.Now().Add(-time.Second)}})

	_, ok := list.Lookup("ip:10.0.0.12")
	assert.False(t, ok)
}

func TestSlotHitByIP(t *testing.T) {
	list := New()
	list.Apply([]core.BanEntry{{Key: "ip:10.0.0.12", Scope: "ip", ExpireAt: time.Now().Add(time.Hour)}})
	s := NewSlot(list)

	d := s.Evaluate(context.Background(), &core.RiskContext{Subject: core.Subject{IP: "10.0.0.12"}})
	assert.Equal(t, core.Ban, d.Verdict)
	assert.Equal(t, []string{"ip"}, d.BanScope)
	assert.Len(t, d.Reasons, 1)
	assert.Equal(t, ReasonID, d.Reasons[0].ID)
}

func TestSlotHitByAttribute(t *testing.T) {
	list := New()
	list.Apply([]core.BanEntry{{Key: "user:u123", Scope: "user", ExpireAt: time.Now().Add(time.Hour)}})
	s := NewSlot(list)

	rc := &core.RiskContext{Subject: core.Subject{Attributes: map[string]string{"user": "u123"}}}
	d := s.Evaluate(context.Background(), rc)
	assert.Equal(t, core.Ban, d.Verdict)
	assert.Equal(t, []string{"user"}, d.BanScope)
}

func TestSlotName(t *testing.T) {
	s := NewSlot(New())
	assert.Equal(t, "banlist", s.Name())
}

func TestSlotMissAllow(t *testing.T) {
	s := NewSlot(New())
	rc := &core.RiskContext{Subject: core.Subject{
		IP:          "10.0.0.99",
		Fingerprint: "f9",
		Attributes:  map[string]string{"user": "u9", "tenant": "t9"},
	}}
	d := s.Evaluate(context.Background(), rc)
	assert.Equal(t, core.Allow, d.Verdict)
}

func TestSubscribeMaintainsList(t *testing.T) {
	n := notifier.NewMemory()
	list := New()
	s := NewSlot(list)

	assert.NoError(t, Subscribe(context.Background(), n, list))

	entries := []core.BanEntry{{Key: "ip:10.0.0.12", Scope: "ip", ExpireAt: time.Now().Add(time.Hour)}}
	payload, _ := json.Marshal(entries)

	assert.NoError(t, n.Publish(context.Background(), notifier.TopicBan, payload))
	rc := &core.RiskContext{Subject: core.Subject{IP: "10.0.0.12"}}
	assert.Equal(t, core.Ban, s.Evaluate(context.Background(), rc).Verdict)

	assert.NoError(t, n.Publish(context.Background(), notifier.TopicUnban, payload))
	assert.Equal(t, core.Allow, s.Evaluate(context.Background(), rc).Verdict)
}
