/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-13 20:08:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-13 20:08:27
 * @FilePath: \go-risk\fingerprint\slot.go
 * @Description: 设备指纹 Slot
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package fingerprint

import (
	"context"

	"github.com/kamalyes/go-risk/core"
)

// slot 设备指纹 Slot：计算指纹与机器人概率，写回请求上下文供后续 Slot 使用，不直接处置
type slot struct {
	fps []Fingerprinter
}

// NewSlot 创建指纹 Slot，聚合多个指纹器（HTTP/TLS 等）
func NewSlot(fps ...Fingerprinter) core.Slot {
	return &slot{fps: fps}
}

// Name 返回 Slot 名称
func (s *slot) Name() string { return "fingerprint" }

// Evaluate 计算指纹并取各指纹器最大的机器人概率，回写上下文后返回观察结论
func (s *slot) Evaluate(ctx context.Context, rc *core.RiskContext) core.Decision {
	var hash string
	var score float64
	for _, fp := range s.fps {
		r := fp.Identify(rc)
		if r.Hash != "" && hash == "" {
			hash = r.Hash
		}
		if r.BotScore > score {
			score = r.BotScore
		}
	}
	if rc.Subject.Fingerprint == "" {
		rc.Subject.Fingerprint = hash
	}
	rc.BotScore = score
	return core.Decision{Verdict: core.Observe}
}