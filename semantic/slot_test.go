/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-02 21:15:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-02 21:15:27
 * @FilePath: \go-risk\semantic\slot_test.go
 * @Description: 语义检测 Slot 单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package semantic

import (
	"context"
	"testing"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestSlotName(t *testing.T) {
	assert.Equal(t, "semantic", NewSlot().Name())
}

func TestSlotEvaluateAllow(t *testing.T) {
	d := NewSlot().Evaluate(context.Background(), &core.RiskContext{})
	assert.Equal(t, core.Allow, d.Verdict)
}

func TestSlotEvaluateBanQuery(t *testing.T) {
	d := NewSlot().Evaluate(context.Background(), &core.RiskContext{Query: "1=1"})
	assert.Equal(t, core.Ban, d.Verdict)
	assert.Equal(t, 95, d.Score)
	assert.Len(t, d.Reasons, 1)
	assert.Equal(t, "semantic-injection", d.Reasons[0].ID)
	assert.Equal(t, "waf", d.Reasons[0].Group)
}

func TestSlotEvaluateBanBody(t *testing.T) {
	d := NewSlot().Evaluate(context.Background(), &core.RiskContext{Body: "<script>"})
	assert.Equal(t, core.Ban, d.Verdict)
}