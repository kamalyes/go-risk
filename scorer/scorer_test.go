/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-20 20:37:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-20 20:37:16
 * @FilePath: \go-risk\scorer\scorer_test.go
 * @Description: 评分引擎单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package scorer

import (
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
)

func TestDecideBan(t *testing.T) {
	cfg := core.DefaultConfig()
	cfg.BanThreshold = 100
	s := New(cfg)
	rc := &core.RiskContext{Subject: core.Subject{IP: "1.2.3.5"}, BotScore: 1.0}
	if d := s.Evaluate(rc); d.Verdict != core.Ban {
		t.Fatalf("verdict = %v, want Ban", d.Verdict)
	}
}

func TestDecideAllow(t *testing.T) {
	s := New(nil)
	rc := &core.RiskContext{Subject: core.Subject{IP: "1.2.3.5"}, BotScore: 0}
	if d := s.Evaluate(rc); d.Verdict != core.Allow {
		t.Fatalf("verdict = %v, want Allow", d.Verdict)
	}
}

func TestDecay(t *testing.T) {
	got := decay(100, time.Hour, 15*time.Minute)
	if got >= 100 {
		t.Fatalf("decay(100, 1h, 15m) = %d, want < 100", got)
	}
	if got <= 0 {
		t.Fatalf("decay = %d, want > 0", got)
	}
}