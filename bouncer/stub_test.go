/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-22 09:37:21
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-22 09:37:21
 * @FilePath: \go-risk\bouncer\stub_test.go
 * @Description: 决策服务单测共享桩（引擎与封禁落地）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bouncer

import (
	"context"

	"github.com/kamalyes/go-risk/core"
)

// stubEngine 固定处置结论的引擎桩
type stubEngine struct {
	verdict core.Decision
}

func (e *stubEngine) Evaluate(context.Context, *core.RiskContext) core.Decision { return e.verdict }
func (e *stubEngine) MarkResult(context.Context, *core.RiskContext, int, bool)  {}
func (e *stubEngine) Close() error                                              { return nil }

// stubSink 记录广播调用的封禁落地桩
type stubSink struct {
	bans   [][]core.BanEntry
	unbans [][]core.BanEntry
}

func (s *stubSink) Ban(_ context.Context, entries []core.BanEntry) error {
	s.bans = append(s.bans, entries)
	return nil
}

func (s *stubSink) Unban(_ context.Context, entries []core.BanEntry) error {
	s.unbans = append(s.unbans, entries)
	return nil
}
