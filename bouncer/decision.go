/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-19 13:36:09
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-19 13:36:09
 * @FilePath: \go-risk\bouncer\decision.go
 * @Description: 处置结论翻译与封禁条目构造
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"net/http"
	"time"

	"github.com/kamalyes/go-risk/core"
)

// allowed 处置结论是否放行，Observe 仅记录不拦截故视为放行
func allowed(d core.Decision) bool {
	return d.Verdict == core.Allow || d.Verdict == core.Observe
}

// denyStatus 拦截结论对应的 HTTP 状态码，挑战与限速返回 429，其余返回 403
func denyStatus(d core.Decision) int {
	if d.Verdict == core.Challenge || d.Verdict == core.Throttle {
		return http.StatusTooManyRequests
	}
	return http.StatusForbidden
}

// banEntries 依封禁维度从风控主体构造封禁条目，ip 与 fingerprint 直接取主体字段
// 其余维度透传扩展属性，缺值时跳过该维度避免提交空键
func banEntries(rc *core.RiskContext, d core.Decision, banDuration time.Duration) []core.BanEntry {
	expire := time.Now().Add(banDuration)
	entries := make([]core.BanEntry, 0, len(d.BanScope))
	for _, scope := range d.BanScope {
		if key := scopeKey(rc, scope); key != "" {
			entries = append(entries, core.BanEntry{Key: key, Scope: scope, ExpireAt: expire})
		}
	}
	return entries
}

// scopeKey 按封禁维度解析主体值，构造 scope:value 形式的封禁键，主体值缺失时返回空
func scopeKey(rc *core.RiskContext, scope string) string {
	switch scope {
	case "ip":
		if rc.Subject.IP != "" {
			return "ip:" + rc.Subject.IP
		}
	case "fingerprint":
		if rc.Subject.Fingerprint != "" {
			return "fingerprint:" + rc.Subject.Fingerprint
		}
	default:
		if v := rc.Subject.Attributes[scope]; v != "" {
			return scope + ":" + v
		}
	}
	return ""
}
