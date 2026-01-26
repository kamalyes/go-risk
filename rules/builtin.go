/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-26 09:33:12
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-26 09:33:12
 * @FilePath: \go-risk\rules\builtin.go
 * @Description: 内置规则集
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import "github.com/kamalyes/go-risk/core"

// Builtin 返回内置规则集（WAF/蜜罐/威胁情报），具体规则在 M1 落地。
func Builtin() core.RuleSnapshot {
	return core.RuleSnapshot{Version: "v0", Rules: []core.Rule{}}
}