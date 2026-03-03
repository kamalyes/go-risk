/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-03 19:23:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-03 19:23:33
 * @FilePath: \go-risk\rules\dynamic.go
 * @Description: 动态规则 Slot，从快照提供方热读取
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import (
	"context"

	"github.com/kamalyes/go-risk/core"
)

// dynamicSlot 每次决策从提供方读取最新快照，支持控制面热更新
type dynamicSlot struct {
	get func() core.RuleSnapshot
}

// NewDynamicSlot 基于快照提供方创建动态 Slot，规则需由提供方保证已排序
func NewDynamicSlot(get func() core.RuleSnapshot) core.Slot {
	return &dynamicSlot{get: get}
}

// Name 返回 Slot 名称
func (s *dynamicSlot) Name() string { return "rules" }

// Evaluate 读取最新快照并返回首个命中规则的处置结论
func (s *dynamicSlot) Evaluate(ctx context.Context, rc *core.RiskContext) core.Decision {
	return evaluate(rc, s.get().Rules)
}
