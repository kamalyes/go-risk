/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-03 21:03:28
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-03 21:03:28
 * @FilePath: \go-risk\engine\engine_test.go
 * @Description: 引擎装配与 Slot 链单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package engine

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/kamalyes/go-risk/core"
	"github.com/kamalyes/go-risk/observe"
	"github.com/stretchr/testify/assert"
)

type mockStore struct{}

func (mockStore) IncrBatch(context.Context, []core.CounterOp) error { return nil }
func (mockStore) WindowCount(context.Context, string, time.Duration) (int64, error) {
	return 0, nil
}
func (mockStore) LoadBanList(context.Context) ([]core.BanEntry, error) { return nil, nil }

type mockNotifier struct{}

func (mockNotifier) Publish(context.Context, string, []byte) error        { return nil }
func (mockNotifier) Subscribe(context.Context, string, func([]byte)) error { return nil }

type mockSlot struct {
	verdict core.Verdict
	score   int
}

func (m mockSlot) Name() string { return "mock" }
func (m mockSlot) Evaluate(context.Context, *core.RiskContext) core.Decision {
	return core.Decision{Verdict: m.verdict, Score: m.score, Reasons: []core.Reason{{ID: "m"}}}
}

func TestNewDefaults(t *testing.T) {
	e := New()
	assert.NotNil(t, e.cfg)
	assert.NotNil(t, e.store)
	assert.NotNil(t, e.notifier)
	assert.NotNil(t, e.chain)
}

func TestWithInjections(t *testing.T) {
	e := New(WithConfig(&core.Config{BanThreshold: 1}), WithStore(mockStore{}), WithNotifier(mockNotifier{}))
	assert.Equal(t, 1, e.cfg.BanThreshold)
	assert.IsType(t, mockStore{}, e.store)
	assert.IsType(t, mockNotifier{}, e.notifier)
}

func TestWithNilOptions(t *testing.T) {
	e := New(WithConfig(nil), WithStore(nil), WithNotifier(nil), WithMetrics(nil), WithLogger(nil))
	assert.Nil(t, e.metrics)
	assert.Nil(t, e.logger)
}

func TestEvaluateWithoutChain(t *testing.T) {
	d := (&Engine{}).Evaluate(context.Background(), &core.RiskContext{})
	assert.Equal(t, core.Allow, d.Verdict)
}

func TestEvaluateWithMetricsLogger(t *testing.T) {
	m := observe.NewMetrics()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	e := New(WithMetrics(m), WithLogger(logger), WithSlots(mockSlot{verdict: core.Ban, score: 50}))
	d := e.Evaluate(context.Background(), &core.RiskContext{})
	assert.Equal(t, core.Ban, d.Verdict)
	assert.Equal(t, int64(1), m.Snapshot(0).Total)
}

func TestWithBuiltinProtection(t *testing.T) {
	e := New(WithBuiltinProtection())
	assert.NotEmpty(t, e.slots)
	d := e.Evaluate(context.Background(), &core.RiskContext{})
	assert.NotEqual(t, core.Ban, d.Verdict)
}

func TestMarkResultAndClose(t *testing.T) {
	e := New()
	e.MarkResult(context.Background(), &core.RiskContext{}, 200, true)
	assert.NoError(t, e.Close())
	assert.NoError(t, (&Engine{}).Close())
}

func TestSlotChainRunShortCircuit(t *testing.T) {
	chain := newSlotChain([]core.Slot{mockSlot{verdict: core.Ban, score: 10}, mockSlot{verdict: core.Allow, score: 99}})
	d := chain.run(context.Background(), &core.RiskContext{})
	assert.Equal(t, core.Ban, d.Verdict)
	assert.Equal(t, 10, d.Score)
	assert.Len(t, d.Reasons, 1)
}

func TestSlotChainRunAggregate(t *testing.T) {
	chain := newSlotChain([]core.Slot{mockSlot{verdict: core.Observe, score: 10}, mockSlot{verdict: core.Observe, score: 20}})
	d := chain.run(context.Background(), &core.RiskContext{})
	assert.Equal(t, core.Allow, d.Verdict)
	assert.Equal(t, 30, d.Score)
	assert.Len(t, d.Reasons, 2)
}