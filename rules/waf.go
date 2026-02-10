/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-10 21:05:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-10 21:05:22
 * @FilePath: \go-risk\rules\waf.go
 * @Description: 内置 WAF 特征规则集
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import "github.com/kamalyes/go-risk/core"

// wafRules 内置 WAF 规则集，覆盖 SQL 注入、XSS、路径遍历、命令注入等常见攻击
// 均声明 urldecode/lower 转换链，命中前先解码，防止编码绕过
func wafRules() []core.Rule {
	return []core.Rule{
		{
			ID: "waf-sql-injection", Priority: 10, Phase: core.PhaseHeaders,
			Targets:    []core.Target{{Collection: "query"}, {Collection: "body"}},
			Transforms: []string{"urldecode", "lower"},
			Op:         core.OpRegex,
			Pattern:    `union\s+select|select\s+.+\s+from|'?\s*or\s+'1'='1'?\s*--|or\s+1\s*=\s*1`,
			Score:      90, Action: core.ActionBan, Group: "waf", Enabled: true,
		},
		{
			ID: "waf-xss", Priority: 20, Phase: core.PhaseHeaders,
			Targets:    []core.Target{{Collection: "query"}, {Collection: "body"}},
			Transforms: []string{"urldecode", "lower"},
			Op:         core.OpRegex,
			Pattern:    `<script|javascript:|onerror=|onload=|onclick=`,
			Score:      85, Action: core.ActionBan, Group: "waf", Enabled: true,
		},
		{
			ID: "waf-path-traversal", Priority: 30, Phase: core.PhaseHeaders,
			Targets:    []core.Target{{Collection: "query"}, {Collection: "body"}, {Collection: "path"}},
			Transforms: []string{"urldecode", "lower"},
			Op:         core.OpRegex,
			Pattern:    `\.\./|\.\.%2f|\.\.\\`,
			Score:      80, Action: core.ActionBan, Group: "waf", Enabled: true,
		},
		{
			ID: "waf-cmd-injection", Priority: 40, Phase: core.PhaseHeaders,
			Targets:    []core.Target{{Collection: "query"}, {Collection: "body"}},
			Transforms: []string{"urldecode", "lower"},
			Op:         core.OpRegex,
			Pattern:    `;wget |;curl |\|cat |\bwhoami\b|\$\(`,
			Score:      88, Action: core.ActionBan, Group: "waf", Enabled: true,
		},
	}
}