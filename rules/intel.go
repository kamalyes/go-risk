/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-11 22:52:39
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-11 22:52:39
 * @FilePath: \go-risk\rules\intel.go
 * @Description: 威胁情报规则
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import (
	"fmt"

	"github.com/kamalyes/go-risk/core"
)

// toolUserAgents 已知攻击工具的 User-Agent 特征，命中即视为自动化攻击
var toolUserAgents = []string{
	"sqlmap", "nmap", "nikto", "dirbuster", "dirsearch",
	"masscan", "zgrab", "nuclei", "acunetix", "wpscan", "hydra",
}

// intelRules 威胁情报规则：命中攻击工具 UA 直接封禁
func intelRules() []core.Rule {
	rs := make([]core.Rule, 0, len(toolUserAgents))
	for i, ua := range toolUserAgents {
		rs = append(rs, core.Rule{
			ID:         fmt.Sprintf("intel-%02d", i+1),
			Priority:   5, Phase: core.PhaseHeaders,
			Targets:    []core.Target{{Collection: "ua"}},
			Transforms: []string{"lower"},
			Op:         core.OpContains,
			Pattern:    ua,
			Score:      100, Action: core.ActionBan, Group: "intel", Enabled: true,
		})
	}
	return rs
}