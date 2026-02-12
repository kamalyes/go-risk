/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-12 21:23:10
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-12 21:23:10
 * @FilePath: \go-risk\rules\slot.go
 * @Description: WAF 规则 Slot
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import (
	"context"
	"sort"

	"github.com/kamalyes/go-risk/core"
)

// slot 基于规则快照的风控 Slot，遍历规则并按优先级返回首个命中的处置
type slot struct {
	snapshot core.RuleSnapshot
}

// NewSlot 基于规则快照创建 Slot，规则按优先级升序排序（数值小者优先）
func NewSlot(snapshot core.RuleSnapshot) core.Slot {
	rules := make([]core.Rule, len(snapshot.Rules))
	copy(rules, snapshot.Rules)
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})
	return &slot{snapshot: core.RuleSnapshot{Version: snapshot.Version, Rules: rules}}
}

// Name 返回 Slot 名称
func (s *slot) Name() string { return "rules" }

// Evaluate 遍历规则并返回首个命中规则的处置结论
func (s *slot) Evaluate(ctx context.Context, rc *core.RiskContext) core.Decision {
	for _, rule := range s.snapshot.Rules {
		hit, matched := Match(rc, rule)
		if !hit {
			continue
		}
		return core.Decision{
			Verdict: verdictOf(rule.Action),
			Score:   rule.Score,
			Reasons: []core.Reason{{ID: rule.ID, Group: rule.Group, Msg: matched}},
		}
	}
	return core.Decision{Verdict: core.Allow}
}

// verdictOf 将规则动作映射为处置分级
func verdictOf(a core.Action) core.Verdict {
	switch a {
	case core.ActionBan:
		return core.Ban
	case core.ActionChallenge:
		return core.Challenge
	case core.ActionThrottle:
		return core.Throttle
	case core.ActionObserve:
		return core.Observe
	default:
		return core.Allow
	}
}
