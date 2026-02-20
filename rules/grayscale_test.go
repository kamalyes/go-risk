/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-20 20:37:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-20 20:37:16
 * @FilePath: \go-risk\rules\grayscale_test.go
 * @Description: 规则灰度单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rules

import (
	"testing"

	"github.com/kamalyes/go-risk/core"
)

func TestGrayHitFull(t *testing.T) {
	rc := &core.RiskContext{Subject: core.Subject{IP: "1.2.3.5"}}
	if !grayHit(rc, 100) {
		t.Fatal("grayHit with pct=100 should be true")
	}
}

func TestGrayHitZero(t *testing.T) {
	rc := &core.RiskContext{Subject: core.Subject{IP: "1.2.3.5"}}
	if !grayHit(rc, 0) {
		t.Fatal("grayHit with pct=0 should be true (full)")
	}
}