/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-27 20:15:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-27 20:15:33
 * @FilePath: \go-risk\control\manager_test.go
 * @Description: 规则快照管理器单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package control

import (
	"testing"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func testRule(id string, prio int) core.Rule {
	return core.Rule{ID: id, Priority: prio, Enabled: true}
}

func TestNewSortsByPriority(t *testing.T) {
	m := New(core.RuleSnapshot{Version: "v1", Rules: []core.Rule{
		testRule("b", 10), testRule("a", 1), testRule("c", 5),
	}})
	cur := m.Current()
	assert.Len(t, cur.Rules, 3)
	assert.Equal(t, "a", cur.Rules[0].ID)
	assert.Equal(t, "c", cur.Rules[1].ID)
	assert.Equal(t, "b", cur.Rules[2].ID)
}

func TestCurrentReturnsCopy(t *testing.T) {
	m := New(core.RuleSnapshot{Rules: []core.Rule{testRule("a", 1)}})
	cur := m.Current()
	cur.Rules[0].ID = "mutated"
	assert.Equal(t, "a", m.Current().Rules[0].ID)
}

func TestReloadAndRollback(t *testing.T) {
	m := New(core.RuleSnapshot{Version: "v1", Rules: []core.Rule{testRule("a", 1)}})
	m.Reload(core.RuleSnapshot{Version: "v2", Rules: []core.Rule{testRule("b", 1)}})
	assert.Equal(t, "b", m.Current().Rules[0].ID)
	assert.True(t, m.Rollback())
	assert.Equal(t, "a", m.Current().Rules[0].ID)
}

func TestRollbackEmptyHistory(t *testing.T) {
	m := New(core.RuleSnapshot{Rules: []core.Rule{testRule("a", 1)}})
	assert.False(t, m.Rollback())
}

func TestReloadUninitialized(t *testing.T) {
	m := &Manager{}
	m.Reload(core.RuleSnapshot{Rules: []core.Rule{testRule("a", 1)}})
	assert.Equal(t, "a", m.Current().Rules[0].ID)
	assert.False(t, m.Rollback())
}