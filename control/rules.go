/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-26 21:32:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-26 21:32:58
 * @FilePath: \go-risk\control\rules.go
 * @Description: 规则增删改、启用与灰度控制
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package control

import "github.com/kamalyes/go-risk/core"

// Upsert 新增或替换规则（按 ID 定位），随后热切换
func (m *Manager) Upsert(rule core.Rule) {
	cur := m.Current()
	rules := append([]core.Rule(nil), cur.Rules...)
	replaced := false
	for i := range rules {
		if rules[i].ID == rule.ID {
			rules[i] = rule
			replaced = true
			break
		}
	}
	if !replaced {
		rules = append(rules, rule)
	}
	m.replace(core.RuleSnapshot{Version: cur.Version, Rules: rules})
}

// Delete 按 ID 删除规则，命中返回 true
func (m *Manager) Delete(id string) bool {
	cur := m.Current()
	idx := -1
	for i := range cur.Rules {
		if cur.Rules[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	rules := cur.Rules
	rules = append(rules[:idx], rules[idx+1:]...)
	m.replace(core.RuleSnapshot{Version: cur.Version, Rules: rules})
	return true
}

// SetEnabled 按 ID 设置规则启用状态，命中返回 true
func (m *Manager) SetEnabled(id string, enabled bool) bool {
	return m.patch(id, func(r *core.Rule) { r.Enabled = enabled })
}

// SetGray 按 ID 设置规则灰度比例，命中返回 true
func (m *Manager) SetGray(id string, pct int) bool {
	return m.patch(id, func(r *core.Rule) { r.GrayPct = pct })
}

// patch 定位规则并原地修改后热切换，命中返回 true
func (m *Manager) patch(id string, fn func(*core.Rule)) bool {
	cur := m.Current()
	rules := append([]core.Rule(nil), cur.Rules...)
	for i := range rules {
		if rules[i].ID == id {
			fn(&rules[i])
			m.replace(core.RuleSnapshot{Version: cur.Version, Rules: rules})
			return true
		}
	}
	return false
}