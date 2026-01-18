/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-18 09:11:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-18 09:11:53
 * @FilePath: \go-risk\core\config.go
 * @Description: 风控引擎配置定义
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

import "time"

// Config 风控引擎配置，阈值分档 + 衰减 + 封禁 + 灰度。
type Config struct {
	ObserveThreshold   int           // 观察阈值（score >= 该值进入观察）
	ThrottleThreshold  int           // 限速阈值
	ChallengeThreshold int           // 挑战阈值
	BanThreshold       int           // 封禁阈值
	DecayHalfLife      time.Duration // 信誉衰减半衰期
	BanDuration        time.Duration // 封禁时长
	BanScope           []string      // 封禁维度：ip/cidr/fingerprint/tenant/user
	GrayPct            int           // 灰度比例 0~100
}

// DefaultConfig 返回默认配置，调用方可按需覆盖字段。
func DefaultConfig() *Config {
	return &Config{
		ObserveThreshold:   30,
		ThrottleThreshold:  60,
		ChallengeThreshold: 80,
		BanThreshold:       100,
		DecayHalfLife:      15 * time.Minute,
		BanDuration:        time.Hour,
		BanScope:           []string{"ip", "fingerprint", "tenant"},
		GrayPct:            100,
	}
}