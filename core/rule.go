/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-16 15:05:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-16 15:05:22
 * @FilePath: \go-risk\core\rule.go
 * @Description: 风控规则、阶段、匹配目标与操作符定义
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

// Phase 规则执行阶段，对齐 Coraza 的分阶段流水线。
type Phase int

const (
	PhaseHeaders Phase = iota // 请求头阶段（先行短路径判定）
	PhaseBody                 // 请求体阶段
	PhaseLogging              // 日志/旁路阶段
)

// Action 规则动作，区分「处置类」与「记录类」。
type Action int

const (
	ActionAllow     Action = iota // 放行
	ActionObserve                 // 记录观察
	ActionThrottle                // 限速
	ActionChallenge               // 挑战
	ActionBan                     // 封禁
)

// Op 匹配操作符。
type Op int

const (
	OpContains   Op = iota // 包含
	OpRegex               // 正则
	OpEquals              // 相等
	OpStartsWith          // 前缀
)

// Target 规则匹配目标，指明命中哪个请求部位。
type Target struct {
	Collection string // 部位：path / query / header / body / arg
	Field      string // 具体字段名（如某个 header 名）
}

// Rule 单条结构化规则，避免重量级表达式依赖；转换链通过名称注入。
type Rule struct {
	ID         string   // 规则 ID
	Priority   int      // 越小越先
	Phase      Phase    // 执行阶段
	Targets    []Target // 匹配目标
	Transforms []string // 命中前解码/规范化转换名称
	Op         Op       // 匹配操作符
	Pattern    string   // 匹配模式
	Score      int      // 命中加分
	Action     Action   // 命中动作
	Group      string   // 规则组（waf/honeypot/intel/behavior）
	Enabled    bool     // 是否启用
	GrayPct    int      // 灰度比例 0~100
}

// RuleSnapshot 规则快照，版本化下发，读多写少无锁切换。
type RuleSnapshot struct {
	Version string
	Rules   []Rule
}