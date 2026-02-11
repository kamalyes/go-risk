/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-11 20:28:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-11 20:28:16
 * @FilePath: \go-risk\rules\honeypot.go
 * @Description: 蜜罐陷阱路径规则
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import (
	"fmt"

	"github.com/kamalyes/go-risk/core"
)

// honeypotPaths 蜜罐陷阱路径，正常用户不会访问，命中即视为扫描器/爬虫
var honeypotPaths = []string{
	"/.git/config",
	"/.env",
	"/.svn/entries",
	"/.htaccess",
	"/.aws/credentials",
	"/phpmyadmin",
	"/wp-login.php",
	"/actuator",
}

// honeypotRules 蜜罐规则：命中陷阱路径直接封禁
func honeypotRules() []core.Rule {
	rs := make([]core.Rule, 0, len(honeypotPaths))
	for i, p := range honeypotPaths {
		rs = append(rs, core.Rule{
			ID:       fmt.Sprintf("honeypot-%02d", i+1),
			Priority: 5, Phase: core.PhaseHeaders,
			Targets: []core.Target{{Collection: "path"}},
			Op:      core.OpStartsWith,
			Pattern: p,
			Score:   100, Action: core.ActionBan, Group: "honeypot", Enabled: true,
		})
	}
	return rs
}