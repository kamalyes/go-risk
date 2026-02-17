/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-17 22:17:39
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-17 22:17:39
 * @FilePath: \go-risk\scorer\slot.go
 * @Description: 分级处置 Slot
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package scorer

import (
	"context"

	"github.com/kamalyes/go-risk/core"
)

// slot 分级处置 Slot：基于评分引擎映射处置分级，位于规则 Slot 之后兜底
type slot struct {
	scorer *Scorer
}

// NewSlot 创建分级处置 Slot
func NewSlot(cfg *core.Config) core.Slot {
	return &slot{scorer: New(cfg)}
}

// Name 返回 Slot 名称
func (s *slot) Name() string { return "scorer" }

// Evaluate 评分并返回处置决策
func (s *slot) Evaluate(ctx context.Context, rc *core.RiskContext) core.Decision {
	return s.scorer.Evaluate(rc)
}
