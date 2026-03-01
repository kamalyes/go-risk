/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-01 15:51:09
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-01 15:51:09
 * @FilePath: \go-risk\semantic\slot.go
 * @Description: 语义检测 Slot
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package semantic

import (
	"context"

	"github.com/kamalyes/go-risk/core"
)

// slot 语义检测 Slot，对请求 query/body 做注入语义分析，命中即封禁
type slot struct {
	det Detector
}

// NewSlot 创建语义检测 Slot
func NewSlot() core.Slot {
	return &slot{det: New()}
}

// Name 返回 Slot 名称
func (s *slot) Name() string { return "semantic" }

// Evaluate 扫描请求并返回处置结论
func (s *slot) Evaluate(ctx context.Context, rc *core.RiskContext) core.Decision {
	for _, v := range []string{rc.Query, rc.Body} {
		if v == "" {
			continue
		}
		if hit, reason := s.det.Scan(v); hit {
			return core.Decision{
				Verdict: core.Ban,
				Score:   95,
				Reasons: []core.Reason{{ID: "semantic-injection", Group: "waf", Msg: reason}},
			}
		}
	}
	return core.Decision{Verdict: core.Allow}
}