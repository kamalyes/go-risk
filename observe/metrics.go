/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-28 10:26:02
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-28 10:26:02
 * @FilePath: \go-risk\observe\metrics.go
 * @Description: 零依赖风控指标
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package observe

import (
	"sort"
	"sync"
	"sync/atomic"

	"github.com/kamalyes/go-risk/core"
)

// RuleHit 规则命中计数
type RuleHit struct {
	ID   string
	Hits int64
}

// Metrics 零依赖内存指标，原子累加，支撑命中率、Verdict 分布与热规则统计
type Metrics struct {
	total     atomic.Int64
	blocked   atomic.Int64
	allow     atomic.Int64
	observe   atomic.Int64
	throttle  atomic.Int64
	challenge atomic.Int64
	ban       atomic.Int64
	mu        sync.Mutex
	rules     map[string]int64
}

// NewMetrics 创建指标
func NewMetrics() *Metrics {
	return &Metrics{rules: make(map[string]int64)}
}

// Record 累计一次决策
func (m *Metrics) Record(d core.Decision) {
	m.total.Add(1)
	switch d.Verdict {
	case core.Allow:
		m.allow.Add(1)
	case core.Observe:
		m.observe.Add(1)
	case core.Throttle:
		m.throttle.Add(1)
	case core.Challenge:
		m.challenge.Add(1)
	case core.Ban:
		m.ban.Add(1)
	}
	if isBlocked(d.Verdict) {
		m.blocked.Add(1)
	}
	m.mu.Lock()
	for _, r := range d.Reasons {
		m.rules[r.ID]++
	}
	m.mu.Unlock()
}

// Snapshot 决策指标快照
type Snapshot struct {
	Total    int64
	Blocked  int64
	Verdicts map[string]int64
	TopRules []RuleHit
}

// Snapshot 返回当前指标快照，TopRules 按命中次数降序取前 topN（topN<=0 全量）
func (m *Metrics) Snapshot(topN int) Snapshot {
	s := Snapshot{
		Total:   m.total.Load(),
		Blocked: m.blocked.Load(),
		Verdicts: map[string]int64{
			"allow":     m.allow.Load(),
			"observe":   m.observe.Load(),
			"throttle":  m.throttle.Load(),
			"challenge": m.challenge.Load(),
			"ban":       m.ban.Load(),
		},
	}
	m.mu.Lock()
	for id, hits := range m.rules {
		s.TopRules = append(s.TopRules, RuleHit{ID: id, Hits: hits})
	}
	m.mu.Unlock()
	sort.Slice(s.TopRules, func(i, j int) bool {
		return s.TopRules[i].Hits > s.TopRules[j].Hits
	})
	if topN > 0 && len(s.TopRules) > topN {
		s.TopRules = s.TopRules[:topN]
	}
	return s
}

// isBlocked 处置类结论（非放行/观察）视为拦截
func isBlocked(v core.Verdict) bool {
	return v != core.Allow && v != core.Observe
}