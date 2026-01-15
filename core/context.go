/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-15 10:37:51
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-15 10:37:51
 * @FilePath: \go-risk\core\context.go
 * @Description: 风控请求主体与上下文契约定义
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

// Subject 风控主体，用于并行维护信誉与行为计数。
// 同一攻击者可能切换 IP/UA，通过多维度主体进行交叉识别。
type Subject struct {
	IP          string // 客户端 IP
	Fingerprint string // 设备指纹（服务端被动指纹或可信客户端指纹）
	TenantID    string // 租户 ID
	UserID      string // 用户 ID
	PlatformID  string // 平台 ID
}

// RiskContext 归一化请求信息（框架无关的纯数据结构）。
// 适配层从各框架请求中提取该结构后交给引擎决策。
type RiskContext struct {
	TraceID     string            // 链路追踪 ID
	Subject     Subject           // 风控主体
	Method      string            // HTTP 方法
	Path        string            // 请求路径
	UserAgent   string            // User-Agent
	Headers     map[string]string // 请求头（去重后小写键）
	HeaderOrder []string          // 请求头出现顺序，供 HTTP 指纹使用
	Query       string            // 原始查询串
	Body        string            // 请求体（M0 原文，后续可替换为摘要）
	BodySize    int64             // 请求体大小
	BotScore    float64           // 预留：设备/机器人概率（0~1），由指纹模块填充
}