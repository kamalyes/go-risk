/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 10:20:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-08 20:15:33
 * @FilePath: \go-risk\fingerprint\fingerprint.go
 * @Description: 设备指纹接口与结果定义
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package fingerprint

import "github.com/kamalyes/go-risk/core"

// Fingerprinter 设备指纹接口，输出稳定指纹与机器人概率。
type Fingerprinter interface {
	// Identify 从请求上下文计算设备指纹与机器人概率。
	Identify(rc *core.RiskContext) Result
}

// Result 指纹计算结果。
type Result struct {
	Hash     string  // 稳定指纹哈希
	BotScore float64 // 机器人概率 0~1
}
