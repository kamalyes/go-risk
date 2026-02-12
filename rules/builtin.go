/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-26 09:33:12
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-12 19:11:53
 * @FilePath: \go-risk\rules\builtin.go
 * @Description: 内置规则集
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import "github.com/kamalyes/go-risk/core"

// Builtin 返回内置规则集（WAF/蜜罐/威胁情报）。
func Builtin() core.RuleSnapshot {
	all := make([]core.Rule, 0, len(wafRules())+len(honeypotRules())+len(intelRules()))
	all = append(all, wafRules()...)
	all = append(all, honeypotRules()...)
	all = append(all, intelRules()...)
	return core.RuleSnapshot{Version: "v1", Rules: all}
}
