/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 09:31:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 09:31:22
 * @FilePath: \go-risk\engine\slotchain.go
 * @Description: Slot 链流水线编排
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package engine

import (
	"context"

	"github.com/kamalyes/go-risk/core"
)

// slotChain 按序执行 Slot 的流水线。
type slotChain struct {
	slots []core.Slot
}

// newSlotChain 创建 Slot 链。
func newSlotChain(slots []core.Slot) *slotChain {
	return &slotChain{slots: slots}
}

// run 执行流水线并聚合决策；处置类结论短路返回。
func (c *slotChain) run(ctx context.Context, rc *core.RiskContext) core.Decision {
	agg := core.Decision{Verdict: core.Allow}
	for _, s := range c.slots {
		d := s.Evaluate(ctx, rc)
		agg.Score += d.Score
		agg.Reasons = append(agg.Reasons, d.Reasons...)
		if d.Verdict != core.Allow && d.Verdict != core.Observe {
			agg.Verdict = d.Verdict
			agg.RetryAfter = d.RetryAfter
			agg.BanScope = d.BanScope
			return agg
		}
	}
	return agg
}