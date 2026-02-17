/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-17 20:51:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-17 20:51:03
 * @FilePath: \go-risk\scorer\scorer.go
 * @Description: 风险评分引擎
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package scorer

import (
	"math"
	"sync"
	"time"

	"github.com/kamalyes/go-risk/core"
)

// Scorer 风险评分引擎：聚合机器人概率与历史信誉，时间衰减后按阈值映射处置分级
type Scorer struct {
	cfg *core.Config
	mu  sync.Mutex
	rep map[string]*reputation
}

type reputation struct {
	score     int
	updatedAt time.Time
}

// New 创建评分引擎
func New(cfg *core.Config) *Scorer {
	if cfg == nil {
		cfg = core.DefaultConfig()
	}
	return &Scorer{cfg: cfg, rep: make(map[string]*reputation)}
}

// Evaluate 对单次请求评分并更新主体信誉，返回处置决策
func (s *Scorer) Evaluate(rc *core.RiskContext) core.Decision {
	key := subjectKey(rc)
	now := time.Now()

	s.mu.Lock()
	rep := s.rep[key]
	if rep == nil {
		rep = &reputation{updatedAt: now}
		s.rep[key] = rep
	}
	rep.score = decay(rep.score, now.Sub(rep.updatedAt), s.cfg.DecayHalfLife)
	rep.updatedAt = now
	rep.score += int(rc.BotScore * 100)
	score := rep.score
	s.mu.Unlock()

	return s.decide(score)
}

// subjectKey 优先使用设备指纹作为主体，其次回退 IP
func subjectKey(rc *core.RiskContext) string {
	if rc.Subject.Fingerprint != "" {
		return "fp:" + rc.Subject.Fingerprint
	}
	return "ip:" + rc.Subject.IP
}

// decide 按阈值分档映射处置分级并填充处置维度
func (s *Scorer) decide(score int) core.Decision {
	d := core.Decision{Score: score}
	switch {
	case score >= s.cfg.BanThreshold:
		d.Verdict = core.Ban
		d.BanScope = s.cfg.BanScope
	case score >= s.cfg.ChallengeThreshold:
		d.Verdict = core.Challenge
	case score >= s.cfg.ThrottleThreshold:
		d.Verdict = core.Throttle
	case score >= s.cfg.ObserveThreshold:
		d.Verdict = core.Observe
	default:
		d.Verdict = core.Allow
	}
	return d
}

// decay 按半衰期衰减信誉分，攻击停止后自动回落解禁
func decay(score int, elapsed, halfLife time.Duration) int {
	if score <= 0 || halfLife <= 0 {
		return score
	}
	factor := math.Pow(0.5, float64(elapsed)/float64(halfLife))
	return int(float64(score) * factor)
}
