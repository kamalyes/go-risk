/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-22 10:15:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-22 10:15:07
 * @FilePath: \go-risk\bouncer\decision_test.go
 * @Description: 处置结论翻译与封禁条目构造单测
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"net/http"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestAllowed(t *testing.T) {
	assert.True(t, allowed(core.Decision{Verdict: core.Allow}))
	assert.True(t, allowed(core.Decision{Verdict: core.Observe}))
	assert.False(t, allowed(core.Decision{Verdict: core.Throttle}))
	assert.False(t, allowed(core.Decision{Verdict: core.Challenge}))
	assert.False(t, allowed(core.Decision{Verdict: core.Ban}))
}

func TestDenyStatus(t *testing.T) {
	assert.Equal(t, http.StatusTooManyRequests, denyStatus(core.Decision{Verdict: core.Throttle}))
	assert.Equal(t, http.StatusTooManyRequests, denyStatus(core.Decision{Verdict: core.Challenge}))
	assert.Equal(t, http.StatusForbidden, denyStatus(core.Decision{Verdict: core.Ban}))
}

func TestScopeKey(t *testing.T) {
	rc := &core.RiskContext{
		Subject: core.Subject{
			IP:          "1.2.3.9",
			Fingerprint: "fp-abc",
			Attributes:  map[string]string{"tenant": "tenant-7"},
		},
	}
	assert.Equal(t, "ip:1.2.3.9", scopeKey(rc, "ip"))
	assert.Equal(t, "fingerprint:fp-abc", scopeKey(rc, "fingerprint"))
	assert.Equal(t, "tenant:tenant-7", scopeKey(rc, "tenant"))
	assert.Empty(t, scopeKey(rc, "user"))
}

func TestBanEntriesSkipsMissingScope(t *testing.T) {
	rc := &core.RiskContext{
		Subject: core.Subject{IP: "1.2.3.9"},
	}
	entries := banEntries(rc, core.Decision{BanScope: []string{"ip", "fingerprint", "tenant"}}, time.Hour)

	assert.Len(t, entries, 1)
	assert.Equal(t, "ip", entries[0].Scope)
	assert.Equal(t, "ip:1.2.3.9", entries[0].Key)
	assert.True(t, entries[0].ExpireAt.After(time.Now()))
}
