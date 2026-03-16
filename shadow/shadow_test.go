/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-16 09:53:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-16 09:53:27
 * @FilePath: \go-risk\shadow\shadow_test.go
 * @Description: 影子模式主备双引擎决策与差异落样单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package shadow

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/require"
)

type mockEngine struct {
	verdict core.Verdict
	score   int

	mu     sync.Mutex
	evals  int
	marks  int
	closed bool
}

func (m *mockEngine) Evaluate(_ context.Context, _ *core.RiskContext) core.Decision {
	m.mu.Lock()
	m.evals++
	m.mu.Unlock()
	return core.Decision{Verdict: m.verdict, Score: m.score}
}

func (m *mockEngine) MarkResult(_ context.Context, _ *core.RiskContext, _ int, _ bool) {
	m.mu.Lock()
	m.marks++
	m.mu.Unlock()
}

func (m *mockEngine) Close() error {
	m.mu.Lock()
	m.closed = true
	m.mu.Unlock()
	return nil
}

func (m *mockEngine) evalCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.evals
}

func TestEvaluateReturnsPrimaryVerdict(t *testing.T) {
	e := New(&mockEngine{verdict: core.Allow}, &mockEngine{verdict: core.Ban})
	d := e.Evaluate(context.Background(), &core.RiskContext{})
	require.Equal(t, core.Allow, d.Verdict)
}

func TestNoSinkEvaluatesOnlyPrimary(t *testing.T) {
	p := &mockEngine{verdict: core.Allow}
	s := &mockEngine{verdict: core.Ban}
	e := New(p, s)
	e.Evaluate(context.Background(), &core.RiskContext{})
	require.Equal(t, 1, p.evalCount())
	require.Equal(t, 0, s.evalCount())
}

func TestDivergentByVerdict(t *testing.T) {
	e := New(&mockEngine{}, &mockEngine{})
	d := e.divergent(
		core.Decision{Verdict: core.Allow, Score: 0},
		core.Decision{Verdict: core.Ban, Score: 0},
	)
	require.True(t, d)
}

func TestDivergentByScoreDelta(t *testing.T) {
	e := New(&mockEngine{}, &mockEngine{}, WithScoreDelta(5))

	require.True(t, e.divergent(
		core.Decision{Verdict: core.Observe, Score: 10},
		core.Decision{Verdict: core.Observe, Score: 20},
	))

	require.False(t, e.divergent(
		core.Decision{Verdict: core.Observe, Score: 10},
		core.Decision{Verdict: core.Observe, Score: 13},
	))
}

func TestDivergentScoreIgnoredWhenDeltaZero(t *testing.T) {
	e := New(&mockEngine{}, &mockEngine{})
	require.False(t, e.divergent(
		core.Decision{Verdict: core.Observe, Score: 0},
		core.Decision{Verdict: core.Observe, Score: 100},
	))
}

func TestSinkReceivesDivergence(t *testing.T) {
	got := make(chan Divergence, 1)
	e := New(
		&mockEngine{verdict: core.Allow},
		&mockEngine{verdict: core.Ban},
		WithSink(func(d Divergence) { got <- d }),
	)
	e.Evaluate(context.Background(), &core.RiskContext{TraceID: "t1", Method: "GET", Path: "/x"})

	select {
	case d := <-got:
		require.Equal(t, core.Allow, d.Primary.Verdict)
		require.Equal(t, core.Ban, d.Shadow.Verdict)
		require.Equal(t, "t1", d.TraceID)
		require.Equal(t, "/x", d.Path)
	case <-time.After(time.Second):
		t.Fatal("差异样本未投递")
	}
	require.NoError(t, e.Close())
}

func TestMarkResultForwardsBoth(t *testing.T) {
	p := &mockEngine{}
	s := &mockEngine{}
	e := New(p, s)
	e.MarkResult(context.Background(), &core.RiskContext{}, 200, true)
	require.Equal(t, 1, p.marks)
	require.Equal(t, 1, s.marks)
}

func TestCloseClosesBoth(t *testing.T) {
	p := &mockEngine{}
	s := &mockEngine{}
	e := New(p, s)
	require.NoError(t, e.Close())
	require.True(t, p.closed)
	require.True(t, s.closed)
}
