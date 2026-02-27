/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-27 21:07:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-27 21:07:52
 * @FilePath: \go-risk\control\rules_test.go
 * @Description: 规则增删改、启用与灰度控制单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package control

import (
	"testing"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestUpsertInsert(t *testing.T) {
	m := New(snapshot())
	m.Upsert(testRule("n", 2))
	assert.Len(t, m.Current().Rules, 2)
}

func TestUpsertReplace(t *testing.T) {
	m := New(snapshot())
	m.Upsert(testRule("a", 99))
	cur := m.Current()
	assert.Len(t, cur.Rules, 1)
	assert.Equal(t, 99, cur.Rules[0].Priority)
}

func TestDeleteHit(t *testing.T) {
	m := New(snapshot())
	assert.True(t, m.Delete("a"))
	assert.Empty(t, m.Current().Rules)
}

func TestDeleteMiss(t *testing.T) {
	m := New(snapshot())
	assert.False(t, m.Delete("missing"))
}

func TestSetEnabledHit(t *testing.T) {
	m := New(snapshot())
	assert.True(t, m.SetEnabled("a", false))
	assert.False(t, m.Current().Rules[0].Enabled)
}

func TestSetEnabledMiss(t *testing.T) {
	m := New(snapshot())
	assert.False(t, m.SetEnabled("missing", false))
}

func TestSetGrayHit(t *testing.T) {
	m := New(snapshot())
	assert.True(t, m.SetGray("a", 25))
	assert.Equal(t, 25, m.Current().Rules[0].GrayPct)
}

func TestSetGrayMiss(t *testing.T) {
	m := New(snapshot())
	assert.False(t, m.SetGray("missing", 25))
}

func snapshot() core.RuleSnapshot {
	return core.RuleSnapshot{Version: "v1", Rules: []core.Rule{testRule("a", 1)}}
}