/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-26 20:11:23
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-26 20:11:23
 * @FilePath: \go-risk\control\manager.go
 * @Description: 规则快照管理器，版本化热更新与回滚
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package control

import (
	"sort"
	"sync"
	"sync/atomic"

	"github.com/kamalyes/go-risk/core"
)

// Manager 规则快照管理器，读路径无锁、写路径串行，控制面与数据面解耦
type Manager struct {
	current atomic.Pointer[core.RuleSnapshot]
	mu      sync.Mutex
	history []core.RuleSnapshot
}

// New 创建管理器并装载初始快照，规则按优先级升序排序
func New(snapshot core.RuleSnapshot) *Manager {
	m := &Manager{}
	cpy := normalize(snapshot)
	m.current.Store(&cpy)
	return m
}

// Current 返回当前快照副本，Slot 热路径可无锁读取
func (m *Manager) Current() core.RuleSnapshot {
	return cloneSnapshot(*m.current.Load())
}

// Reload 热切换为新快照，旧快照压入历史栈用于回滚
func (m *Manager) Reload(snapshot core.RuleSnapshot) {
	m.replace(snapshot)
}

// Rollback 回滚到最近一个历史快照，无历史返回 false
func (m *Manager) Rollback() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := len(m.history)
	if n == 0 {
		return false
	}
	prev := m.history[n-1]
	m.history = m.history[:n-1]
	cpy := cloneSnapshot(prev)
	m.current.Store(&cpy)
	return true
}

// replace 将当前快照压栈并切换为新快照
func (m *Manager) replace(snapshot core.RuleSnapshot) {
	cpy := normalize(snapshot)
	m.mu.Lock()
	defer m.mu.Unlock()
	if cur := m.current.Load(); cur != nil {
		m.history = append(m.history, cloneSnapshot(*cur))
	}
	m.current.Store(&cpy)
}

// normalize 深拷贝规则并按优先级升序排序
func normalize(s core.RuleSnapshot) core.RuleSnapshot {
	rules := append([]core.Rule(nil), s.Rules...)
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})
	s.Rules = rules
	return s
}

// cloneSnapshot 深拷贝快照，避免外部修改影响内部状态
func cloneSnapshot(s core.RuleSnapshot) core.RuleSnapshot {
	s.Rules = append([]core.Rule(nil), s.Rules...)
	return s
}