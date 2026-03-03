/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-03 20:01:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-03 20:01:18
 * @FilePath: \go-risk\rules\dynamic_test.go
 * @Description: 动态规则 Slot 单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import (
	"context"
	"testing"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestDynamicSlotName(t *testing.T) {
	s := NewDynamicSlot(func() core.RuleSnapshot { return core.RuleSnapshot{} })
	assert.Equal(t, "rules", s.Name())
}

func TestDynamicSlotEvaluateHit(t *testing.T) {
	s := NewDynamicSlot(func() core.RuleSnapshot {
		return core.RuleSnapshot{Rules: []core.Rule{{
			ID: "d", Priority: 1, Phase: core.PhaseHeaders,
			Targets: []core.Target{{Collection: "path"}},
			Op:      core.OpEquals, Pattern: "/bad", Score: 100, Action: core.ActionBan, Enabled: true,
		}}}
	})
	d := s.Evaluate(context.Background(), &core.RiskContext{Path: "/bad"})
	assert.Equal(t, core.Ban, d.Verdict)
}

func TestDynamicSlotEvaluateMiss(t *testing.T) {
	s := NewDynamicSlot(func() core.RuleSnapshot { return core.RuleSnapshot{} })
	d := s.Evaluate(context.Background(), &core.RiskContext{Path: "/bad"})
	assert.Equal(t, core.Allow, d.Verdict)
}
