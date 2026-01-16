/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-16 09:22:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-16 09:22:08
 * @FilePath: \go-risk\core\decision.go
 * @Description: 风控处置分级、决策结果与计数三态定义
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

import "time"

// Verdict 处置分级（递进式响应），对齐业界观察→限速→挑战→封禁的渐进策略。
type Verdict int

const (
	Allow Verdict = iota // 放行
	Observe              // 观察（埋点记录，不拦截）
	Throttle             // 限速（动态收紧）
	Challenge            // 挑战（429 / 验证码 / 二次验证）
	Ban                  // 封禁
)

func (v Verdict) String() string {
	switch v {
	case Allow:
		return "allow"
	case Observe:
		return "observe"
	case Throttle:
		return "throttle"
	case Challenge:
		return "challenge"
	case Ban:
		return "ban"
	}
	return "allow"
}

// CountVerict 计数三态结果，用于分布式聚合区分「临界」与「超限」。
type CountVerict int

const (
	CountUnknown  CountVerict = iota // 后端异常
	CountAllowed                     // 允许
	CountEdge                        // 临界（本窗口最后一单）
	CountRejected                    // 拒绝
)

// Reason 命中原因，供审计与日志输出。
type Reason struct {
	ID    string // 命中的规则 ID
	Group string // 规则组（waf/honeypot/intel/behavior）
	Msg   string // 命中说明
}

// Decision 决策结果，等价于一次风控中断指令；Verdict=Allow 表示无需动作。
type Decision struct {
	Verdict    Verdict        // 分级结论
	Score      int            // 风险评分
	Reasons    []Reason       // 命中原因链
	RetryAfter time.Duration  // 挑战/限速的等待时间建议
	BanScope   []string       // 封禁维度（ip/cidr/fingerprint/tenant/user）
}