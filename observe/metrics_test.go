/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-28 16:23:35
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-28 16:23:35
 * @FilePath: \go-risk\observe\metrics_test.go
 * @Description: 风控指标单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package observe

import (
	"testing"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestMetricsRecordVerdicts(t *testing.T) {
	m := NewMetrics()
	m.Record(core.Decision{Verdict: core.Allow})
	m.Record(core.Decision{Verdict: core.Observe})
	m.Record(core.Decision{Verdict: core.Throttle})
	m.Record(core.Decision{Verdict: core.Challenge})
	m.Record(core.Decision{Verdict: core.Ban})

	s := m.Snapshot(0)
	assert.Equal(t, int64(5), s.Total)
	assert.Equal(t, int64(3), s.Blocked)
	for k, want := range map[string]int64{"allow": 1, "observe": 1, "throttle": 1, "challenge": 1, "ban": 1} {
		assert.Equal(t, want, s.Verdicts[k], "verdict %s", k)
	}
}

func TestMetricsRecordReasons(t *testing.T) {
	m := NewMetrics()
	m.Record(core.Decision{Verdict: core.Ban, Reasons: []core.Reason{{ID: "r1"}, {ID: "r2"}}})
	m.Record(core.Decision{Verdict: core.Ban, Reasons: []core.Reason{{ID: "r1"}}})

	s := m.Snapshot(0)
	assert.Len(t, s.TopRules, 2)
	assert.Equal(t, RuleHit{ID: "r1", Hits: 2}, s.TopRules[0])
	assert.Equal(t, RuleHit{ID: "r2", Hits: 1}, s.TopRules[1])
}

func TestMetricsSnapshotTopN(t *testing.T) {
	m := NewMetrics()
	for i := 0; i < 3; i++ {
		m.Record(core.Decision{Verdict: core.Ban, Reasons: []core.Reason{{ID: "a"}}})
	}
	m.Record(core.Decision{Verdict: core.Ban, Reasons: []core.Reason{{ID: "b"}}})

	s := m.Snapshot(1)
	assert.Len(t, s.TopRules, 1)
	assert.Equal(t, "a", s.TopRules[0].ID)
}

func TestIsBlocked(t *testing.T) {
	assert.False(t, isBlocked(core.Allow))
	assert.False(t, isBlocked(core.Observe))
	assert.True(t, isBlocked(core.Throttle))
	assert.True(t, isBlocked(core.Challenge))
	assert.True(t, isBlocked(core.Ban))
}