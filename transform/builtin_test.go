/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-29 10:37:19
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-29 10:37:19
 * @FilePath: \go-risk\transform\builtin_test.go
 * @Description: 内置转换单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package transform

import "testing"

func TestLower(t *testing.T) {
	if got := Lower()("AbC"); got != "abc" {
		t.Fatalf("Lower = %q, want abc", got)
	}
}

func TestURLDecode(t *testing.T) {
	if got := URLDecode()("%2Fetc%2Fpasswd"); got != "/etc/passwd" {
		t.Fatalf("URLDecode = %q, want /etc/passwd", got)
	}
}

func TestNormalizePath(t *testing.T) {
	if got := NormalizePath()("/a/../b"); got != "/b" {
		t.Fatalf("NormalizePath = %q, want /b", got)
	}
}

func TestChainApply(t *testing.T) {
	c := NewChain(Lower(), NormalizePath())
	if got := c.Apply("/A/../B"); got != "/b" {
		t.Fatalf("Chain.Apply = %q, want /b", got)
	}
}