/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-02 20:33:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-02 20:33:16
 * @FilePath: \go-risk\semantic\engine_test.go
 * @Description: 语义检测引擎单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package semantic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanEmpty(t *testing.T) {
	hit, reason := New().Scan("")
	assert.False(t, hit)
	assert.Equal(t, "", reason)
}

func TestScanUnescapeError(t *testing.T) {
	hit, _ := New().Scan("%")
	assert.False(t, hit)
}

func TestScanSQLComment(t *testing.T) {
	for _, in := range []string{"a--b", "a/*x*/", "a#b"} {
		r, ok := scanSQL(in)
		assert.True(t, ok, in)
		assert.Equal(t, "sql-comment", r, in)
	}
}

func TestScanSQLKeyword(t *testing.T) {
	r, ok := scanSQL("select name from t")
	assert.True(t, ok)
	assert.Equal(t, "sql-keyword", r)
}

func TestScanSQLTautology(t *testing.T) {
	r, ok := scanSQL("1=1")
	assert.True(t, ok)
	assert.Equal(t, "sql-tautology", r)
}

func TestScanSQLNoHit(t *testing.T) {
	r, ok := scanSQL("hello")
	assert.False(t, ok)
	assert.Equal(t, "", r)
}

func TestScanXSSProtocol(t *testing.T) {
	r, ok := scanXSS("javascript:alert(1)")
	assert.True(t, ok)
	assert.Equal(t, "xss-protocol", r)
}

func TestScanXSSTag(t *testing.T) {
	r, ok := scanXSS("<script>")
	assert.True(t, ok)
	assert.Equal(t, "xss-tag", r)
}

func TestScanXSSHandler(t *testing.T) {
	r, ok := scanXSS("onerror=")
	assert.True(t, ok)
	assert.Equal(t, "xss-handler", r)
}

func TestScanXSSNoHit(t *testing.T) {
	r, ok := scanXSS("hello")
	assert.False(t, ok)
	assert.Equal(t, "", r)
}

func TestScanCommandInjection(t *testing.T) {
	r, ok := scanCommand("; ls")
	assert.True(t, ok)
	assert.Equal(t, "cmd-injection", r)
}

func TestScanCommandNoSep(t *testing.T) {
	r, ok := scanCommand("hello")
	assert.False(t, ok)
	assert.Equal(t, "", r)
}

func TestScanCommandSepNoExec(t *testing.T) {
	r, ok := scanCommand("a;b")
	assert.False(t, ok)
	assert.Equal(t, "", r)
}

func TestScanIntegrated(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"' or '1'='1", true},
		{"<img src=x onerror=alert(1)>", true},
		{"ping -c 1 127.0.0.1 && whoami", true},
		{"hello world", false},
	}
	for _, c := range cases {
		hit, _ := New().Scan(c.in)
		assert.Equal(t, c.want, hit, c.in)
	}
}

func TestTautology(t *testing.T) {
	assert.True(t, tautology("1=1"))
	assert.True(t, tautology("'a'='a'"))
	assert.False(t, tautology("hello"))
}

func TestQuoteTautology(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"'a'='a'", true},
		{"'a'='b'", false},
		{`"x"="x"`, true},
		{"xyz", false},
		{"'abc", false},
		{"''=''", false},
		{"'a' = 'a'", true},
		{`'a'="a"`, false},
		{"'a'x", false},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, quoteTautology(c.in), c.in)
	}
}

func TestTokenContains(t *testing.T) {
	assert.True(t, tokenContains("select name from t", "select"))
	assert.False(t, tokenContains("select name from t", "where"))
	assert.False(t, tokenContains("xselect", "select"))
}

func TestIsNumeric(t *testing.T) {
	assert.False(t, isNumeric(""))
	assert.True(t, isNumeric("123"))
	assert.False(t, isNumeric("12a"))
}

func TestTokensKeepEquals(t *testing.T) {
	assert.Equal(t, []string{"a", "=", "1", "b"}, tokens("a=1 b"))
}
