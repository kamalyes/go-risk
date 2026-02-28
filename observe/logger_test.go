/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-28 17:02:38
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-28 17:02:38
 * @FilePath: \go-risk\observe\logger_test.go
 * @Description: 单行 KV 结构化决策日志单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package observe

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/kamalyes/go-risk/core"
	"github.com/stretchr/testify/assert"
)

func TestLogNilLogger(t *testing.T) {
	Log(context.Background(), nil, &core.RiskContext{}, core.Decision{})
}

func TestLogOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	rc := &core.RiskContext{TraceID: "t1", Subject: core.Subject{Fingerprint: "fp", IP: "1.2.3.5"}}
	Log(context.Background(), logger, rc, core.Decision{
		Verdict: core.Ban,
		Score:   95,
		Reasons: []core.Reason{{ID: "r1"}, {ID: "r2"}},
	})
	out := buf.String()
	for _, want := range []string{"t1", "fp", "ban", "r1,r2"} {
		assert.Contains(t, out, want)
	}
}

func TestSubjectFingerprint(t *testing.T) {
	rc := &core.RiskContext{Subject: core.Subject{Fingerprint: "fp", IP: "1.2.3.5"}}
	assert.Equal(t, "fp", subject(rc))
}

func TestSubjectIP(t *testing.T) {
	rc := &core.RiskContext{Subject: core.Subject{IP: "1.2.3.5"}}
	assert.Equal(t, "1.2.3.5", subject(rc))
}

func TestReasonsEmpty(t *testing.T) {
	assert.Equal(t, "", reasons(nil))
}

func TestReasonsJoin(t *testing.T) {
	assert.Equal(t, "a,b", reasons([]core.Reason{{ID: "a"}, {ID: "b"}}))
}