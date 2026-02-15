/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-15 20:15:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-15 20:15:33
 * @FilePath: \go-risk\rules\matcher_test.go
 * @Description: 规则匹配单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import (
	"testing"

	"github.com/kamalyes/go-risk/core"
)

func TestWafDetectsSQLInjection(t *testing.T) {
	rc := &core.RiskContext{Query: "id=1 union select password from users"}
	for _, r := range wafRules() {
		if hit, _ := Match(rc, r); hit {
			return
		}
	}
	t.Fatal("expected SQL injection to be detected")
}

func TestWafDetectsUrlEncodedXSS(t *testing.T) {
	rc := &core.RiskContext{Query: "q=%3Cscript%3Ealert(1)%3C%2Fscript%3E"}
	for _, r := range wafRules() {
		if r.ID == "waf-xss" {
			if hit, _ := Match(rc, r); !hit {
				t.Fatal("expected encoded XSS to be detected after urldecode")
			}
		}
	}
}

func TestHoneypotPath(t *testing.T) {
	rc := &core.RiskContext{Path: "/.git/config"}
	for _, r := range honeypotRules() {
		if hit, _ := Match(rc, r); hit {
			return
		}
	}
	t.Fatal("expected honeypot path to be detected")
}

func TestSlotBanDecision(t *testing.T) {
	s := NewSlot(core.RuleSnapshot{Rules: []core.Rule{{
		ID: "t", Priority: 1, Phase: core.PhaseHeaders,
		Targets: []core.Target{{Collection: "path"}},
		Op:      core.OpEquals, Pattern: "/bad", Score: 100, Action: core.ActionBan, Enabled: true,
	}}})
	d := s.Evaluate(nil, &core.RiskContext{Path: "/bad"})
	if d.Verdict != core.Ban {
		t.Fatalf("verdict = %v, want Ban", d.Verdict)
	}
}