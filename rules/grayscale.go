/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-18 19:25:38
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-18 19:25:38
 * @FilePath: \go-risk\rules\grayscale.go
 * @Description: 规则灰度判定
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import "github.com/kamalyes/go-risk/core"

// grayHit 按主体灰度判定规则是否对当前请求生效
// pct 为 0 或 100 表示全量，否则按主体哈希取模判定
func grayHit(rc *core.RiskContext, pct int) bool {
	if pct <= 0 || pct >= 100 {
		return true
	}
	key := rc.Subject.Fingerprint
	if key == "" {
		key = rc.Subject.IP
	}
	return fnv32(key)%100 < uint32(pct)
}

func fnv32(key string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}
	return h
}
