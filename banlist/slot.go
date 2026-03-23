/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-23 10:27:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-23 10:27:52
 * @FilePath: \go-risk\banlist\slot.go
 * @Description: 封禁名单拦截 Slot，链首对名单内主体确定性拒绝
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package banlist

import (
	"context"

	"github.com/kamalyes/go-risk/core"
)

// slot 封禁名单拦截 Slot，命中即返回封禁，短路后续规则与评分
type slot struct {
	list *List
}

// NewSlot 基于封禁名单创建拦截 Slot
func NewSlot(list *List) core.Slot {
	return &slot{list: list}
}

// Name 返回 Slot 名称
func (s *slot) Name() string { return "banlist" }

// Evaluate 按主体维度构造候选键查询名单，命中即返回封禁
func (s *slot) Evaluate(_ context.Context, rc *core.RiskContext) core.Decision {
	for _, key := range candidateKeys(rc) {
		if e, ok := s.list.Lookup(key); ok {
			return core.Decision{
				Verdict:  core.Ban,
				Reasons:  []core.Reason{{ID: ReasonID, Group: ReasonGroup, Msg: key}},
				BanScope: []string{e.Scope},
			}
		}
	}
	return core.Decision{Verdict: core.Allow}
}

// candidateKeys 由风控主体构造封禁候选键，覆盖 ip 指纹与扩展属性维度
func candidateKeys(rc *core.RiskContext) []string {
	keys := make([]string, 0, 2+len(rc.Subject.Attributes))
	if rc.Subject.IP != "" {
		keys = append(keys, "ip:"+rc.Subject.IP)
	}
	if rc.Subject.Fingerprint != "" {
		keys = append(keys, "fingerprint:"+rc.Subject.Fingerprint)
	}
	for k, v := range rc.Subject.Attributes {
		if v != "" {
			keys = append(keys, k+":"+v)
		}
	}
	return keys
}
