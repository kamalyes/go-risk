/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-15 20:15:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-15 20:15:33
 * @FilePath: \go-risk\fingerprint\http_test.go
 * @Description: HTTP 指纹单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package fingerprint

import (
	"testing"

	"github.com/kamalyes/go-risk/core"
)

func TestHTTPIdentifiesEmptyUA(t *testing.T) {
	fp := NewHTTP()
	rc := &core.RiskContext{UserAgent: "", Headers: map[string]string{}}
	r := fp.Identify(rc)
	if r.BotScore <= 0 {
		t.Fatalf("BotScore = %f, want > 0 for empty UA", r.BotScore)
	}
	if r.Hash == "" {
		t.Fatal("Hash should not be empty")
	}
}

func TestHTTPIdentifiesBotUA(t *testing.T) {
	fp := NewHTTP()
	rc := &core.RiskContext{UserAgent: "sqlmap/1.0", Headers: map[string]string{"Accept-Language": "zh"}}
	r := fp.Identify(rc)
	if r.BotScore < 0.3 {
		t.Fatalf("BotScore = %f, want >= 0.3 for bot UA", r.BotScore)
	}
}

func TestHTTPStableHash(t *testing.T) {
	fp := NewHTTP()
	rc := &core.RiskContext{
		UserAgent:   "Mozilla/5.0",
		HeaderOrder: []string{"Host", "User-Agent"},
		Headers:     map[string]string{"Accept-Language": "zh", "Accept-Encoding": "gzip"},
	}
	a := fp.Identify(rc).Hash
	b := fp.Identify(rc).Hash
	if a != b {
		t.Fatalf("hash unstable: %s vs %s", a, b)
	}
}