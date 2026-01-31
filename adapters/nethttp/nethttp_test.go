/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-31 09:55:31
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-31 09:55:31
 * @FilePath: \go-risk\adapters\nethttp\nethttp_test.go
 * @Description: net/http 适配器闭环单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package nethttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kamalyes/go-risk/engine"
)

func TestWrapPassthrough(t *testing.T) {
	e := engine.New()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := Wrap(e, next)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestExtractBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader("hello"))
	rc := extract(req)
	if rc.Body != "hello" {
		t.Fatalf("body = %q, want hello", rc.Body)
	}
}